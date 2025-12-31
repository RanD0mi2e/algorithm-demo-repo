package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const proxyPort = "7895"

func main() {
	addr := ":" + proxyPort
	log.Printf("Go HTTPS (CONNECT) 代理服务器启动，监听地址: %s", addr)

	// 使用 http.Server 来处理请求，特别是 CONNECT 方法
	server := &http.Server{
		Addr: addr,
		// 使用一个专门处理 CONNECT 请求的 Handler
		Handler: http.HandlerFunc(handleConnect),
		// 读写超时设置
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// 处理所有请求，重点是 CONNECT 方法
func handleConnect(w http.ResponseWriter, r *http.Request) {
	log.Printf("收到请求: %s %s 来自: %s", r.Method, r.URL.Host, r.RemoteAddr)

	// 1. 检查是否是 CONNECT 方法（HTTPS 代理的关键）
	if r.Method == http.MethodConnect {
		// 目标地址格式为 host:port，例如 example.com:443
		targetAddr := r.URL.Host

		// 确保目标地址有效
		if !strings.Contains(targetAddr, ":") {
			targetAddr += ":443" // CONNECT 默认是 HTTPS 端口
		}

		// 2. 先 Hijack 客户端连接
		// 必须在建立目标连接前 Hijack，这样可以完全控制响应的发送
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			log.Println("ResponseWriter 不支持 Hijacker")
			http.Error(w, "无法创建隧道", http.StatusInternalServerError)
			return
		}
		clientConn, _, err := hijacker.Hijack()
		if err != nil {
			log.Printf("Hijack 失败: %v", err)
			return
		}
		defer clientConn.Close()

		// 3. 建立与目标服务器的 TCP 连接
		destConn, err := net.DialTimeout("tcp", targetAddr, 10*time.Second)
		if err != nil {
			// 手动发送错误响应
			clientConn.Write([]byte("HTTP/1.1 503 Service Unavailable\r\n\r\n"))
			log.Printf("连接失败到 %s: %v", targetAddr, err)
			return
		}
		defer destConn.Close()

		// 4. 手动发送 200 Connection Established 响应
		// 这是 CONNECT 方法的标准响应格式
		_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
		if err != nil {
			log.Printf("发送响应失败: %v", err)
			return
		}

		// 5. 启动双向数据转发 (隧道)
		// 使用 channel 来同步两个方向的转发完成
		done := make(chan struct{}, 2)
		go transferWithSignal(destConn, clientConn, done)
		go transferWithSignal(clientConn, destConn, done)

		// 等待其中一个方向完成（通常是其中一方关闭连接）
		<-done
		// 此时另一个方向也会因为连接关闭而结束

	} else {
		// 对于非 CONNECT 请求（普通 HTTP），可以进行简单的 HTTP 转发
		handleHTTP(w, r)
	}
}

// 双向数据转发函数（带完成信号）
func transferWithSignal(destination io.Writer, source io.Reader, done chan struct{}) {
	// 将数据从 source 复制到 destination
	// io.Copy 会一直阻塞直到 EOF 或发生错误
	io.Copy(destination, source)

	// 转发完成，发送信号
	// 使用 select 防止向已关闭的 channel 发送
	select {
	case done <- struct{}{}:
	default:
	}
}

// 简单的 HTTP 转发函数 (与前一个 MVP 类似)
func handleHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("处理 HTTP 请求: %s %s", r.Method, r.URL)

	r.RequestURI = ""
	r.Header.Del("Proxy-Connection")

	client := &http.Client{}
	resp, err := client.Do(r)

	if err != nil {
		http.Error(w, "HTTP 代理转发失败: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 复制响应头和体
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

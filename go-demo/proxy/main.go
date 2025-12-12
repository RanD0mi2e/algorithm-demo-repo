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
			targetAddr += ":80" // 默认 HTTP 端口
		}

		// 2. 建立与目标服务器的 TCP 连接
		destConn, err := net.DialTimeout("tcp", targetAddr, 5*time.Second)
		if err != nil {
			http.Error(w, "无法连接到目标服务器: "+err.Error(), http.StatusServiceUnavailable)
			log.Printf("连接失败到 %s: %v", targetAddr, err)
			return
		}

		// 3. 告知客户端连接已建立
		// 代理回复 200 OK，表示隧道已准备好
		w.WriteHeader(http.StatusOK)

		// 4. 劫持客户端连接 (Hijack)
		// 这一步是为了拿到原始的底层 TCP 连接，而不是 HTTP 响应流
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			log.Println("ResponseWriter 不支持 Hijacker")
			http.Error(w, "无法创建隧道", http.StatusInternalServerError)
			return
		}
		clientConn, _, err := hijacker.Hijack()
		if err != nil {
			log.Printf("Hijack 失败: %v", err)
			http.Error(w, "无法劫持客户端连接", http.StatusInternalServerError)
			return
		}

		// 5. 启动双向数据转发 (隧道)
		// 一旦隧道建立，代理的工作就是将字节从一端复制到另一端，不再关心 HTTP 协议。
		go transfer(destConn, clientConn)
		go transfer(clientConn, destConn)

	} else {
		// 对于非 CONNECT 请求（普通 HTTP），可以进行简单的 HTTP 转发
		handleHTTP(w, r)
	}
}

// 双向数据转发函数
func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()

	// 将数据从 source 复制到 destination
	// io.Copy 会一直阻塞直到 EOF 或发生错误
	if _, err := io.Copy(destination, source); err != nil {
		// 在隧道关闭时，通常会收到 io.EOF 或 net.OpError，可以忽略
		if !strings.Contains(err.Error(), "use of closed network connection") {
			// log.Printf("数据转发错误: %v", err) // 生产环境可以根据需要开启
		}
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

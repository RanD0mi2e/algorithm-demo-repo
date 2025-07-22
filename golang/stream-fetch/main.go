package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func simulateLLMAPI(prompt string, resultChan chan<- string, doneChan chan<- bool) {
	// 模拟响应内容
	response := "这是对问题'" + prompt + "'的回答：\n\n"
	response += "流式文本 输出效果是通过服务器发送事件(SSE)或WebSocket实现的。这种技术允许服务器向客户端推送数据，而不需要客户端发起多次请求。在实际应用中，每个字符或词组会被单独发送，从而创造出打字机效果。"
	words := strings.Split(response, "")
	for _, word := range words {

		resultChan <- word + " "
		time.Sleep(50 * time.Millisecond)
	}
	doneChan <- true
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	r.GET("/stream-fetch", func(c *gin.Context) {
		prompt := c.Query("prompt")
		if prompt == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缺少prompt参数"})
			return
		}
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.WriteHeader(http.StatusOK)

		// 创建通道
		resultChan := make(chan string, 10) // 缓冲区大小可以根据需要调整
		doneChan := make(chan bool)
		go simulateLLMAPI(prompt, resultChan, doneChan)

		for {
			select {
			case text := <-resultChan:
				fmt.Fprintf(c.Writer, "data: %s\n\n", text)
				c.Writer.Flush()
			case <-doneChan:
				// 发送完成事件
				fmt.Fprintf(c.Writer, "event: done\ndata: {\"done\": true}\n\n")
				c.Writer.Flush()
				close(resultChan)
				close(doneChan)
				return
			case <-c.Request.Context().Done():
				return
			}
		}
	})

	if err := r.Run(":30002"); err != nil {
		fmt.Println("服务启动失败:", err)
	}
}

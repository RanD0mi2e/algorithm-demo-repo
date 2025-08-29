package main

import (
	"fmt"
	hwt "go-demo/hash-wheel-timer"
	"time"
)

func main() {
	start := time.Now()
	exec := hwt.NewExecutor(4, 60)
	timer := hwt.NewHashedWheelTimer(1*time.Second, 8, exec)
	timer.Start()
	// 短任务（3s）
	timer.NewTimeout(3*time.Second, func() {
		fmt.Printf("time elipse1: %v\n", time.Since(start))
		fmt.Println("[短任务] at", time.Now())
	})
	time.Sleep(1 * time.Second)
	// 长任务（20s，跨轮）
	timer.NewTimeout(20*time.Second, func() {
		fmt.Printf("time elipse2: %v\n", time.Since(start))
		fmt.Println("[跨轮长任务] at", time.Now())
	})
	time.Sleep(25 * time.Second)
	timer.Stop()
}

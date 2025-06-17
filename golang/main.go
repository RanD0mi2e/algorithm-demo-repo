package main

import (
	"algorithm/fibonacci"
	"fmt"
)

func main() {
	mtx := fibonacci.Matrix{A: 1, B: 1, C: 1, D: 0}
	res := fibonacci.GetNthFibonacciNumByMatrix(5, mtx)
	fmt.Println(res)
}

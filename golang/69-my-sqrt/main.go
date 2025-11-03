package main

import (
	"fmt"
	"math"
)

func main() {
	input := 8
	output := mySqrt(input)
	fmt.Printf("result: %v", output)
}

// 牛顿迭代法
/**
以直代曲思想：
1.获取二次曲线x0=c（c为常数）做一条斜率为2x0的曲线
2.f'(xi) = 2xi
3.求斜率直线的与x轴交点，逐渐逼近曲线根
y = f'(xi)(x - xi) + (xi^2 - c)
4.把坐标(x0, x0^2-c)带入直线方程,可得到
2xi(x-xi) + (xi^2 - c) = 0
=> x = 1/2 * (xi + c/xi)
5.误差精度只要 > 10 ^ -7 即可
*/

func mySqrt(x int) int {
	if x == 0 {
		return x
	}
	c, x0 := float64(x), float64(x)
	for {
		x1 := 0.5 * (x0 + c/x0)
		if math.Abs(x0-x1) < 1e-7 {
			break
		}
		x0 = x1
	}
	return int(x0)
}

// 二分法
func mySqrt2(x int) int {
	if x < 0 {
		return -1
	}
	if x < 2 {
		return x
	}
	left, right := 0, x
	for left < right {
		mid := left + (right-left)/2
		square := mid * mid
		if square > x {
			right = mid
		} else if square < x {
			left = mid + 1
		} else {
			return mid
		}
	}
	return left - 1
}

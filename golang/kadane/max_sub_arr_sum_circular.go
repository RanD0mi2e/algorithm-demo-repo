package main

import (
	"fmt"
	"math"
)

func main() {
	input := []int{1, -2, 3, -2}
	output := maxSubArrSumCircular(input)
	fmt.Printf("result: %d", output)
}

func maxSubArrSumCircular(nums []int) int {
	// 数组长度为0，默认最大子数组和为0
	n := len(nums)
	if n == 0 {
		return 0
	}

	// 环形数组看成两倍nums长度的线性数组，并求整个线性数组的前缀和
	var prefix = make([]int, 2*n+1)
	for i := 0; i < 2*n; i++ {
		prefix[i+1] = prefix[i] + nums[i%n]
	}

	// 假如存在[i,j]范围最大子数组和
	// 环形子数组最大和 = prefix[j+1] - prefix[i]
	// for循环判断每个元素作为j+1的位置，此时j+1对应的前缀值固定，要使得子数组和最大，则prefix[i]要尽可能小
	// 为了确保prefix[i]尽可能小，而且子数组的范围不能超过nums长度
	// 原条件可以变成：在2n数组长度范围内，以循环中每个j+1为右界，找到滑动窗口为n的范围内最小的前缀和索引i
	// 为了保证找到这个最小的前缀和索引i，可以利用单调队列单调递增特性，确保队列中左界队首为前缀和最小的值
	var result int = math.MinInt
	var deque = make([]int, 0)
	deque = append(deque, 0)
	for i := 1; i <= 2*n; i++ {
		if len(deque) > 0 {
			// 单调队列左界队首索引已经超出了滑动窗口左界，这时候这个值不能作为可用值，必须出队
			if deque[0] < i-n {
				deque = deque[1:]
			}
			// 比较每次 prefix[j+1] - prefix[i] 的差值，维护最大值
			curPrefixSum := prefix[i] - prefix[deque[0]]
			result = max(curPrefixSum, result)
			// 为了保证队列单调性，把所有大于当前有界的前缀和全部出列
			for prefix[i] <= prefix[len(deque)-1] {
				deque = deque[:len(deque)-1]
			}
		}
		deque = append(deque, i)
	}

	return result
}

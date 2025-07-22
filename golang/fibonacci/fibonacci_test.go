package fibonacci

import (
	"testing"
)

func TestGetNthFibonacciNum(t *testing.T) {
	tests := []struct {
		n        int
		expected int
	}{
		{0, 0},   // 边界条件：n 为 0
		{1, 1},   // 边界条件：n 为 1
		{2, 1},   // 边界条件：n 为 2
		{3, 2},   // 边界条件：n 为 3
		{10, 55}, // 一般情况：n 为 10
	}

	for _, test := range tests {
		result := GetNthFibonacciNum(test.n)
		if result != test.expected {
			t.Errorf("GetNthFibonacciNum(%d) = %d; expected %d", test.n, result, test.expected)
		}
	}
}

package primefactor

import (
	"reflect"
	"testing"
)

func TestGetPrimeFactorSequence(t *testing.T) {
	tests := []struct {
		input    int
		expected []int
	}{
		{1, []int{}},                    // 没有素因子
		{2, []int{2}},                   // 唯一的素数
		{3, []int{3}},                   // 素数
		{4, []int{2, 2}},                // 合成数
		{6, []int{2, 3}},                // 合成数
		{8, []int{2, 2, 2}},             // 合成数
		{9, []int{3, 3}},                // 合成数
		{10, []int{2, 5}},               // 合成数
		{11, []int{11}},                 // 素数
		{100, []int{2, 2, 5, 5}},        // 大合成数
		{1000, []int{2, 2, 2, 5, 5, 5}}, // 大合成数
		{-10, []int{}},                  // 负数（假设返回空切片）
		{0, []int{}},                    // 0（假设返回空切片）
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := GetPrimeFactorSequence(test.input)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("For input %d, expected %v, but got %v", test.input, test.expected, result)
			}
		})
	}
}

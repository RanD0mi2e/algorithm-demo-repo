package LCM

import "testing"

func TestLowestCommonMultiple(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{12, 18, 36},
		{5, 10, 10},
		{7, 3, 21},
		{21, 6, 42},
		{8, 9, 72},
	}

	for _, test := range tests {
		result := LowestCommonMultiple(test.a, test.b)
		if result != test.expected {
			t.Errorf("LowestCommonMultiple(%d, %d) = %d; expected %d", test.a, test.b, result, test.expected)
		}
	}
}

func TestLowestCommonMultipleForNums(t *testing.T) {
	tests := []struct {
		nums     []int
		expected int
	}{
		{[]int{12, 18, 24}, 72},
		{[]int{5, 10, 15}, 30},
		{[]int{7, 3, 14}, 42},
		{[]int{21, 6, 14}, 42},
		{[]int{8, 9, 21}, 504},
	}

	for _, test := range tests {
		result := LowestCommonMultipleForNums(test.nums...)
		if result != test.expected {
			t.Errorf("LowestCommonMultipleForNums(%v) = %d; expected %d", test.nums, result, test.expected)
		}
	}
}

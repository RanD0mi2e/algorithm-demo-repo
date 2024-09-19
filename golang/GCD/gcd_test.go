package GCD

import (
	"testing"
)

func TestGreatestCommonDivisor(t *testing.T) {
	tests := []struct {
		name           string
		a, b, expected int
	}{
		{"48 and 18", 48, 18, 6},
		{"56 and 98", 56, 98, 14},
		{"101 and 103", 101, 103, 1},
		{"0 and 5", 0, 5, 5},
		{"5 and 0", 5, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GreatestCommonDivisor(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("greatestCommonDivisor(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestGreatestCommonDivisorForNums(t *testing.T) {
	tests := []struct {
		nums     []int
		expected int
	}{
		{[]int{12, 15, 21}, 3},
		{[]int{100, 25, 50}, 25},
		{[]int{7, 14, 21}, 7},
		{[]int{9, 27, 81}, 9},
		{[]int{17, 19, 23}, 1},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := GreatestCommonDivisorForNums(test.nums...)
			if result != test.expected {
				t.Errorf("GreatestCommonDivisorForNums(%v) = %d; want %d", test.nums, result, test.expected)
			}
		})
	}
}

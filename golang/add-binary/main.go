package main

import "fmt"

func main() {
	testCase := [][]string{
		{"0", "1"},
		{"1", "1"},
		{"0", "0"},
		{"1011", "1101"},
	}

	for index, tc := range testCase {
		result := addBinary(tc[0], tc[1])
		fmt.Printf("case %d result: %v\n", index+1, result)
	}
}

func addBinary(a string, b string) string {
	result := ""
	carry := 0
	i, j := len(a)-1, len(b)-1
	for i >= 0 || j >= 0 || carry > 0 {
		sum := carry
		if i >= 0 {
			sum += int(a[i] - '0')
			i--
		}
		if j >= 0 {
			sum += int(b[j] - '0')
			j--
		}
		result = fmt.Sprintf("%d", sum%2) + result
		carry = sum / 2
	}

	return result
}

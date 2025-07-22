package main

import "fmt"

func findMaxFrequentNum(str string) byte {
	counts := [10]int{}
	var mostFrequent byte
	maxCount := 0

	for _, count := range str {
		if count >= '0' && count <= '9' {
			digit := int(count - '0')
			counts[digit]++

			if counts[digit] > maxCount {
				maxCount = counts[digit]
				mostFrequent = byte(count)
			}
		}
	}

	return mostFrequent
}

func findModeNum(str string) string {
	var candidate byte
	count := 0
	for _, c := range str {
		if c < '0' && c > '9' {
			continue
		}
		if count == 0 {
			candidate = byte(c)
			count++
		} else if candidate == byte(c) {
			count++
		} else {
			count--
		}
	}

	if checkIsOverHalf(candidate, str) {
		return string(candidate)
	} else {
		return "-1"
	}
}

func checkIsOverHalf(char byte, str string) bool {
	count := 0
	for _, c := range str {
		if char != byte(c) {
			continue
		}
		count++
	}
	return count > len(str)/2
}

func main() {
	input := "123115411"
	result := findMaxFrequentNum(input)
	fmt.Printf("出现最多的数字是: %c\n", result)

	result2 := findModeNum(input)
	fmt.Printf("众数是: %s", result2)
}

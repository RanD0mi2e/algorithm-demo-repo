package main

import "fmt"

func letterCombinations(digits string) []string {
	if len(digits) == 0 {
		return []string{}
	}

	phoneMap := map[byte]string{
		'2': "abc",
		'3': "def",
		'4': "ghi",
		'5': "jkl",
		'6': "mno",
		'7': "pqrs",
		'8': "tuv",
		'9': "wxyz",
	}

	result := []string{}

	var backTrace func(index int, path string)
	backTrace = func(index int, path string) {
		if index == len(digits) {
			result = append(result, path)
			return
		}

		digit := digits[index]
		letters := phoneMap[digit]

		for i := range letters {
			backTrace(index+1, path+string(letters[i]))
		}
	}

	backTrace(0, "")

	return result
}

func main() {
	fmt.Println("示例1:", letterCombinations("23"))
	fmt.Println("示例2:", letterCombinations(""))
	fmt.Println("示例3:", letterCombinations("2"))
}

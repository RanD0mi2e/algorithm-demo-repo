package binarysearch

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func ReadFileToSlice(filename string) ([]int, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var numberList []int

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, ",")

		for _, item := range items {
			item = strings.TrimSpace(item)
			num, err := strconv.Atoi(item)

			if err != nil {
				return nil, err
			}

			numberList = append(numberList, num)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return numberList, nil
}

func Rank(arr []int, target int) int {

	left, right := 0, len(arr)-1
	for left <= right {
		mid := left + (right-left)/2
		if target < arr[mid] {
			right = mid - 1
		} else if target > arr[mid] {
			left = mid + 1
		} else {
			return mid
		}
	}

	return -1
}

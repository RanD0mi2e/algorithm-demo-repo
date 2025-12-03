package main

import (
	"fmt"
	"sort"
)

func main() {
	input := []string{"eat", "tea", "tan", "ate", "nat", "bat"}
	output := groupAnagrams2(input)
	fmt.Println("result:", output)
}

func groupAnagrams(strs []string) [][]string {
	repeatMap := make(map[string]int)
	var result [][]string
	for _, item := range strs {
		newItem := stringSort(item)
		if _, ok := repeatMap[newItem]; !ok {
			result = append(result, []string{item})
			repeatMap[newItem] = len(result) - 1
		} else {
			idx := repeatMap[newItem]
			result[idx] = append(result[idx], item)
		}
	}
	return result
}

func stringSort(str string) string {
	r := []rune(str)
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return string(r)
}

func groupAnagrams2(strs []string) [][]string {
	anagrams := make(map[[26]int][]string)

	for _, str := range strs {
		var tempMap [26]int
		for _, char := range str {
			tempMap[char-'a']++
		}
		anagrams[tempMap] = append(anagrams[tempMap], str)
	}

	result := [][]string{}
	for _, group := range anagrams {
		result = append(result, group)
	}

	return result
}

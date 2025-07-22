package main

func findSubstring(s string, words []string) []int {
	result := []int{}

	if len(s) == 0 || len(words) == 0 {
		return result
	}

	wordCount := len(words)
	wordLen := len(words[0])
	totalLen := wordCount * wordLen

	if len(s) < totalLen {
		return result
	}

	wordFreq := make(map[string]int)
	for _, word := range words {
		wordFreq[word]++
	}

	for start := 0; start < wordLen; start++ {
		left, right := start, start
		currWordFreq := make(map[string]int)
		count := 0

		for right+wordLen <= len(s) {
			word := s[right : right+wordLen]
			right += wordLen

			if freq, exists := wordFreq[word]; exists {
				currWordFreq[word]++
				if currWordFreq[word] <= freq {
					count++
				} else {
					for currWordFreq[word] > freq {
						leftWord := s[left : left+wordLen]
						currWordFreq[leftWord]--
						left += wordLen

						if currWordFreq[leftWord] < wordFreq[leftWord] {
							count--
						}
					}
				}

				if count == wordCount {
					result = append(result, left)
					leftWord := s[left : left+wordLen]
					currWordFreq[leftWord]--
					count--
					left += wordLen
				}
			} else {
				currWordFreq = make(map[string]int)
				count = 0
				left = right
			}
		}
	}

	return result
}

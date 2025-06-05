package main

func permute(nums []int) [][]int {
	if len(nums) == 0 {
		return [][]int{}
	}

	var result [][]int
	var current []int
	used := make([]bool, len(nums))

	var dfs func()
	dfs = func() {
		if len(nums) == len(current) {
			temp := make([]int, len(nums))
			copy(temp, current)
			result = append(result, temp)
			return
		}

		for i, num := range nums {
			if !used[i] {
				current = append(current, num)
				used[i] = true
				dfs()
				used[i] = false
				current = current[:len(current)-1]
			}
		}
	}

	dfs()

	return result
}

func main() {}

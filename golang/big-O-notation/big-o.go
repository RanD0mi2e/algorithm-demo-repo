package main

import "fmt"

func main() {
	list := []int{1, 2, 3}
	permute(list, 0, len(list)-1)
}

func permute(list []int, l, r int) {
	if l == r {
		for i := 0; i <= r; i++ {
			fmt.Printf("%v ", list[i])
		}
		fmt.Printf("\n")
	} else {
		for i := l; i <= r; i++ {
			// swap
			list[l], list[i] = list[i], list[l]
			permute(list, l+1, r)
			list[l], list[i] = list[i], list[l]
		}
	}
}

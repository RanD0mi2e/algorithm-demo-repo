package main

import (
	"algorithm/LCM"
	"fmt"
)

func main() {

	res := LCM.LowestCommonMultipleForNums(2, 6, 8)
	fmt.Printf("result: %d", res)
}

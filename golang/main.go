package main

import (
	primefactor "algorithm/prime-factor"
	"fmt"
)

func main() {

	res := primefactor.GetPrimeFactorSequence(9)
	fmt.Println(res)
}

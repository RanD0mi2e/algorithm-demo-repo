package primefactor

func GetPrimeFactorSequence(n int) []int {
	result := []int{}

	// 素数范围 > 0
	if n <= 0 {
		return result
	}

	// 从最小素数2去整除
	for n%2 == 0 {
		result = append(result, 2)
		n = n / 2
	}

	// 3 ~ +MAX 整除素数， 直到 i > sqrt(n)
	for i := 3; i*i <= n; i += 2 {
		for n%i == 0 {
			result = append(result, i)
			n = n / i
		}
	}

	if n > 2 {
		result = append(result, n)
	}

	return result
}

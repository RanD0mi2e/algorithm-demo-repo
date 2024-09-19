package GCD

// GCD(a,b) = GCD(b, a%b) when b = 0 return a

// 两数的最大公约数
func GreatestCommonDivisor(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func GreatestCommonDivisorForNums(nums ...int) int {
	length := len(nums)
	result := nums[0]
	for i := 1; i < length; i++ {
		result = GreatestCommonDivisor(result, nums[i])
	}
	return result
}

func GreatestCommonDivisor2(a, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	result := min(a, b)
	for result > 0 {
		if a%result == 0 && b%result == 0 {
			break
		}
		result--
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

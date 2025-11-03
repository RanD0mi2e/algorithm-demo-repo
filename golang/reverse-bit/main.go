package reversebit

func reverseBits(n int) int {
	result := 0
	for range 32 {
		result = (result << 1) | (n & 1)
		n = n >> 1
	}
	return result
}

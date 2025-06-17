package fibonacci

// 获取序号n的斐波那契数
func GetNthFibonacciNum(n int) int {
	lp, rp := 0, 1
	if n == 0 {
		return lp
	}
	if n == 1 {
		return rp
	}
	temp := -1
	for i := 2; i <= n; i++ {
		temp = lp + rp
		lp, rp = rp, temp
	}
	return temp
}

type Matrix struct {
	A, B, C, D int
}

// 获取序号n的斐波那契数列（线代解法）
func GetNthFibonacciNumByMatrix(n int, mtx Matrix) Matrix {
	res := Matrix{
		1, 0, 0, 1,
	}

	base := mtx

	for n > 0 {
		if n%2 == 1 {
			res = matrixMultiple(res, base)
		}
		base = matrixMultiple(base, base)
		n /= 2
	}

	return res
}

// 矩阵乘法
func matrixMultiple(x, y Matrix) Matrix {
	return Matrix{
		x.A*y.A + x.B*y.C,
		x.A*y.B + x.B*y.D,
		x.C*y.A + x.D*y.C,
		x.C*y.B + x.D*y.D,
	}
}

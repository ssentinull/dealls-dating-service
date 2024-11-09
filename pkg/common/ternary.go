package common

func Ternary[T any](cond bool, valTrue, valFalse T) T {
	if cond {
		return valTrue
	}
	return valFalse
}

func Fallback[T comparable](x T, y T) T {
	zero := *new(T)

	if x == zero {
		return y
	}

	return x
}

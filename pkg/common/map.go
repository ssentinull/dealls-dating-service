package common

func CleanMap[K comparable, V comparable](m map[K]V) map[K]V {
	res := make(map[K]V)

	for key, value := range m {
		var zero V

		if value == zero {
			continue
		}

		res[key] = value
	}

	return res
}

func MapToSlice[K comparable, V any](m map[K]V) []V {
	res := make([]V, 0)

	for _, value := range m {
		res = append(res, value)
	}

	return res
}

package common

import (
	"strconv"
)

func ToValue[T interface{}](v *T) T {
	if v == nil {
		v = new(T)
	}

	return *v
}

func IntToBool(v int) bool {
	return v != 0
}

func StringToBool(s string) bool {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return v
}

func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return int64(0)
	}
	return i
}

func StringToFloat64(v string) float64 {
	float, _ := strconv.ParseFloat(v, 64)
	return float
}

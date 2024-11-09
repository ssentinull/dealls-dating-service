package common

import (
	"fmt"
	"reflect"
)

func GetIndexWithFieldValue[T any](arr []T, fieldName string, targetValue interface{}) (int, error) {
	for i, elem := range arr {
		elemValue := reflect.ValueOf(elem)
		fieldValue := elemValue.FieldByName(fieldName)
		if fieldValue.IsValid() && fieldValue.Type().AssignableTo(reflect.TypeOf(targetValue)) && fieldValue.Interface() == targetValue {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Target value not found in slice")
}

func RemoveAtIndex[T any](slice []T, s int) ([]T, error) {
	if slice == nil {
		return nil, fmt.Errorf("input slice is nil")
	}
	if s < 0 || s >= len(slice) {
		return nil, fmt.Errorf("index out of range")
	}
	return append(slice[:s], slice[s+1:]...), nil
}

func SliceHas[T comparable](slice []T, x T) bool {
	if slice == nil {
		return false
	}

	for _, element := range slice {
		if element == x {
			return true
		}
	}

	return false
}

func CleanSlice[T comparable](slice []T) []T {
	res := make([]T, 0)

	zero := *new(T)

	for _, s := range slice {
		if s == zero {
			continue
		}

		if SliceHas(res, s) {
			continue
		}

		res = append(res, s)
	}

	return res
}

func MakeSliceFromFunc[T any](N int, f func(i int) T) []T {
	res := make([]T, N)

	for i := 0; i < N; i++ {
		res[i] = f(i)
	}

	return res
}

func FirstOrZero[T any](arr []T) T {
	if len(arr) > 0 {
		return arr[0]
	}

	zero := *new(T)

	return zero
}

func SliceToMap[K comparable, V any](slice []V, f func(V) K) map[K]V {
	res := make(map[K]V, 0)

	for _, s := range slice {
		key := f(s)

		res[key] = s
	}

	return res
}

func SliceFilter[S any](slice []S, f func(src S) bool) []S {
	res := make([]S, 0)

	for _, s := range slice {
		if f(s) {
			res = append(res, s)
		}
	}

	return res
}

func UnpackSlice[S any](slice []S, dest ...*S) {
	for i, v := range slice {
		if i == len(dest) {
			break
		}

		*dest[i] = v
	}
}

func SliceTransform[S any, D any](slice []S, f func(src S) (dest D)) []D {
	res := make([]D, 0)

	for _, s := range slice {
		res = append(res, f(s))
	}

	return res
}

func SliceAdd[T comparable](slice []T, x T) []T {
	if !SliceHas(slice, x) {
		slice = append(slice, x)
	}

	return slice
}

func SliceRemove[T comparable](slice []T, x T) []T {
	return SliceFilter(slice, func(src T) bool {
		return src != x
	})
}

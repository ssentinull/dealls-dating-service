package common

func ToPointer[T any](x T) *T {
	return &x
}

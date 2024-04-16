package utils

func Ternary[T any](truth bool, truthy, falsy T) T {
	if truth {
		return truthy
	}

	return falsy
}

package utils

func Map[T, R any](arr []T, delegate func(row T, i int) R) []R {
	res := make([]R, len(arr))
	for i := 0; i < len(arr); i++ {
		res[i] = delegate(arr[i], i)
	}

	return res
}

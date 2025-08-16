package utils

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func RepeatSlice[T any](n int, def T) []T {
	list := make([]T, n)
	for i := range list {
		list[i] = def
	}
	return list
}

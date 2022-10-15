package pkg

func Type2pointer[T comparable](val T) *T {
	return &val
}

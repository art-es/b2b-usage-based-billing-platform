package ptr

func To[T any](v T) *T {
	return &v
}

func Value[T any](v *T) T {
	if v != nil {
		return *v
	}

	var defval T
	return defval
}

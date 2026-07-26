package ptr

func To[T any](v T) *T {
	return &v
}

func ToOrNil[T comparable](v T) *T {
	var defval T
	if v == defval {
		return nil
	}

	return &v
}

func Value[T any](v *T) T {
	if v != nil {
		return *v
	}

	var defval T
	return defval
}

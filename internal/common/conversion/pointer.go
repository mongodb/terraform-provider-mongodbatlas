package conversion

func Pointer[T any](x T) *T {
	return &x
}

func StringPtr(v string) *string {
	if v != "" {
		return &v
	}
	return nil
}

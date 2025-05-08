package conversion

func Pointer[T any](x T) *T {
	return &x
}

func IntPtr(v int) *int {
	if v != 0 {
		return &v
	}
	return nil
}

func StringPtr(v string) *string {
	if v != "" {
		return &v
	}
	return nil
}

func SliceFromPtr[T any](slicePtr *[]T) []T {
	if slicePtr == nil {
		return []T{}
	}
	return *slicePtr
}

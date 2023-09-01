package util

import "time"

// utility conversions that can potentially be defined in sdk
func TimePtrToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	res := t.String()
	return &res
}

func Int64PtrToIntPtr(i64 *int64) *int {
	if i64 == nil {
		return nil
	}

	i := int(*i64)
	return &i
}

package util

import "time"

// utility conversions that can potentially be defined in sdk
func TimePtrToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	res := TimeToString(*t)
	return &res
}

// TimeToString returns a RFC3339 date time string format.
// The resulting format is identical to the format returned by Atlas API, documented as ISO 8601 timestamp format in UTC.
// Example format: "2023-07-18T16:12:23Z"
func TimeToString(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func Int64PtrToIntPtr(i64 *int64) *int {
	if i64 == nil {
		return nil
	}

	i := int(*i64)
	return &i
}

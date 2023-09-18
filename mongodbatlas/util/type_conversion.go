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
// It also returns decimals in seconds (up to nanoseconds) if available.
// Example formats: "2023-07-18T16:12:23Z", "2023-07-18T16:12:23.456Z"
func TimeToString(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

func Int64PtrToIntPtr(i64 *int64) *int {
	if i64 == nil {
		return nil
	}

	i := int(*i64)
	return &i
}

// IsStringPresent returns true if the string is non-empty.
func IsStringPresent(strPtr *string) bool {
	return strPtr != nil && len(*strPtr) > 0
}

package conversion

import (
	"strings"
	"time"
)

func SafeString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

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

func IntPtrToInt64Ptr(i *int) *int64 {
	if i == nil {
		return nil
	}

	i64 := int64(*i)
	return &i64
}

// IsStringPresent returns true if the string is non-empty.
func IsStringPresent(strPtr *string) bool {
	return strPtr != nil && *strPtr != ""
}

// MongoDBRegionToAWSRegion converts region in US_EAST_1-like format to us-east-1-like
func MongoDBRegionToAWSRegion(region string) string {
	return strings.ReplaceAll(strings.ToLower(region), "_", "-")
}

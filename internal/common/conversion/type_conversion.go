package conversion

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func SafeValue[T any](v *T) T {
	if v != nil {
		return *v
	}
	var emptyValue T
	return emptyValue
}

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

// StringToTime is the opposite to TimeToString, returns ok if conversion is possible.
func StringToTime(str string) (time.Time, bool) {
	ret, err := time.Parse(time.RFC3339Nano, str)
	return ret, err == nil
}

func StringPtrToTimePtr(str *string) (*time.Time, bool) {
	if str == nil {
		return nil, true
	}
	res, ok := StringToTime(*str)
	return &res, ok
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

// AWSRegionToMongoDBRegion converts region in us-east-1-like format to US_EAST_1-like
func AWSRegionToMongoDBRegion(region string) string {
	return strings.ReplaceAll(strings.ToUpper(region), "-", "_")
}

type TFPrimitiveType interface {
	IsUnknown() bool
}

func NilForUnknown[T any](primitiveAttr TFPrimitiveType, value *T) *T {
	if primitiveAttr.IsUnknown() {
		return nil
	}
	return value
}

func NilForUnknownOrEmptyString(primitiveAttr types.String) *string {
	value := NilForUnknown(primitiveAttr, primitiveAttr.ValueStringPointer())
	if value == nil || *value == "" {
		return nil
	}
	return value
}

func TFModelList[T any](ctx context.Context, diags *diag.Diagnostics, input types.List) []T {
	elements := make([]T, len(input.Elements()))
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	return elements
}

func TFModelObject[T any](ctx context.Context, diags *diag.Diagnostics, input types.Object) *T {
	item := new(T)
	if localDiags := input.As(ctx, item, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	return item
}

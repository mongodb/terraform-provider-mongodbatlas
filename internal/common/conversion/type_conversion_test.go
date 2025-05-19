package conversion_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func TestTimeWithoutNanos(t *testing.T) {
	inputTime := time.Date(2023, time.July, 18, 16, 12, 23, 0, time.UTC)
	expectedOutput := "2023-07-18T16:12:23Z"

	result := conversion.TimeToString(inputTime)
	assert.Equal(t, expectedOutput, result)

	expectedTime, ok := conversion.StringToTime(result)
	assert.True(t, ok)
	assert.Equal(t, expectedTime, inputTime)
}

func TestTimeWithNanos(t *testing.T) {
	inputTime := time.Date(2023, time.July, 18, 16, 12, 23, 456_000_000, time.UTC)
	expectedOutput := "2023-07-18T16:12:23.456Z"

	result := conversion.TimeToString(inputTime)
	assert.Equal(t, expectedOutput, result)

	expectedTime, ok := conversion.StringToTime(result)
	assert.True(t, ok)
	assert.Equal(t, expectedTime, inputTime)
}

func TestStringToTimeInvalid(t *testing.T) {
	_, ok := conversion.StringToTime("")
	assert.False(t, ok)

	_, ok = conversion.StringToTime("123")
	assert.False(t, ok)
}

func TestIsStringPresent(t *testing.T) {
	var (
		empty    = ""
		oneBlank = " "
		str      = "text"
	)
	tests := []struct {
		strPtr   *string
		expected bool
	}{
		{nil, false},
		{&empty, false},
		{&oneBlank, true},
		{&str, true},
	}
	for _, test := range tests {
		if resp := conversion.IsStringPresent(test.strPtr); resp != test.expected {
			t.Errorf("IsStringPresent(%v) = %v; want %v", test.strPtr, resp, test.expected)
		}
	}
}

func TestMongoDBRegionToAWSRegion(t *testing.T) {
	tests := []struct {
		region   string
		expected string
	}{
		{"US_EAST_1", "us-east-1"},
		{"us-east-1", "us-east-1"},
		{"nothing", "nothing"},
	}

	for _, test := range tests {
		if resp := conversion.MongoDBRegionToAWSRegion(test.region); resp != test.expected {
			t.Errorf("MongoDBRegionToAWSRegion(%v) = %v; want %v", test.region, resp, test.expected)
		}
	}
}

func TestAWSRegionToMongoDBRegion(t *testing.T) {
	tests := []struct {
		region   string
		expected string
	}{
		{"us-east-1", "US_EAST_1"},
		{"US-EAST-1", "US_EAST_1"},
	}

	for _, test := range tests {
		if resp := conversion.AWSRegionToMongoDBRegion(test.region); resp != test.expected {
			t.Errorf("AWSRegionToMongoDBRegion(%v) = %v; want %v", test.region, resp, test.expected)
		}
	}
}

func TestSafeValue(t *testing.T) {
	var boolPointer *bool
	assert.False(t, conversion.SafeValue(boolPointer))
	trueBool := true
	assert.True(t, conversion.SafeValue(&trueBool))
	var intPointer *int
	assert.Equal(t, 0, conversion.SafeValue(intPointer))
	var stringPointer *string
	assert.Empty(t, conversion.SafeValue(stringPointer))
}

func TestNilForUnknownOrEmpty(t *testing.T) {
	assert.Nil(t, conversion.NilForUnknownOrEmptyString(types.StringPointerValue(nil)))
	emptyString := ""
	assert.Nil(t, conversion.NilForUnknownOrEmptyString(types.StringPointerValue(&emptyString)))
	testString := "test"
	assert.Equal(t, "test", *conversion.NilForUnknownOrEmptyString(types.StringPointerValue(&testString)))
}

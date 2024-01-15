package conversion_test

import (
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
)

func TestTimeToStringWithoutNanos(t *testing.T) {
	inputTime := time.Date(2023, time.July, 18, 16, 12, 23, 0, time.UTC)
	expectedOutput := "2023-07-18T16:12:23Z"

	result := conversion.TimeToString(inputTime)

	if result != expectedOutput {
		t.Errorf("TimeToString(%v) = %v; want %v", inputTime, result, expectedOutput)
	}
}

func TestTimeToStringWithNanos(t *testing.T) {
	inputTime := time.Date(2023, time.July, 18, 16, 12, 23, 456_000_000, time.UTC)
	expectedOutput := "2023-07-18T16:12:23.456Z"

	result := conversion.TimeToString(inputTime)

	if result != expectedOutput {
		t.Errorf("TimeToString(%v) = %v; want %v", inputTime, result, expectedOutput)
	}
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

func TestNonEmptySliceToSlicePtr(t *testing.T) {
	var (
		nilSlice         []string
		emptyNonNilSlice = []string{}
	)
	tests := []struct {
		name     string
		expected *[]string
		given    []string
	}{
		{"nil pointer", nil, nil},
		{"nil slice", nil, nilSlice},
		{"empty non-nil slice", nil, emptyNonNilSlice},
		{"slice with content", &[]string{"hello", "there"}, []string{"hello", "there"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, conversion.NonEmptySliceToSlicePtr(tc.given))
		})
	}
}

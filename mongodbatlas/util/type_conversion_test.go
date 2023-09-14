package util_test

import (
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/util"
)

func TestTimeToStringWithoutNanos(t *testing.T) {
	inputTime := time.Date(2023, time.July, 18, 16, 12, 23, 0, time.UTC)
	expectedOutput := "2023-07-18T16:12:23Z"

	result := util.TimeToString(inputTime)

	if result != expectedOutput {
		t.Errorf("TimeToString(%v) = %v; want %v", inputTime, result, expectedOutput)
	}
}

func TestTimeToStringWithNanos(t *testing.T) {
	inputTime := time.Date(2023, time.July, 18, 16, 12, 23, 456_000_000, time.UTC)
	expectedOutput := "2023-07-18T16:12:23.456Z"

	result := util.TimeToString(inputTime)

	if result != expectedOutput {
		t.Errorf("TimeToString(%v) = %v; want %v", inputTime, result, expectedOutput)
	}
}

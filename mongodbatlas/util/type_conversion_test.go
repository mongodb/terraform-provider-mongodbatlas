package util_test

import (
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/util"
)

func TestTimeToString(t *testing.T) {
	inputTime := time.Date(2021, time.September, 12, 15, 4, 5, 0, time.UTC)
	expectedOutput := "2021-09-12T15:04:05Z"

	result := util.TimeToString(inputTime)

	if result != expectedOutput {
		t.Errorf("TimeToString(%v) = %v; want %v", inputTime, result, expectedOutput)
	}
}

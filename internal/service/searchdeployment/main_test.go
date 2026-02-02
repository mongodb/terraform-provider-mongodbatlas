package searchdeployment_test

import (
	"os"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestMain(m *testing.M) {
	os.Exit(acc.Run(m))
}

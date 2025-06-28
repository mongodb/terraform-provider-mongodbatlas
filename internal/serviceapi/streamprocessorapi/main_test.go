package streamprocessorapi_test

import (
	"os"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestMain(m *testing.M) {
	acc.SetupSharedResources()
	exitCode := m.Run()
	os.Exit(exitCode)
}

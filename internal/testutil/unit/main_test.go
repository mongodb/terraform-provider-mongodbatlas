package unit_test

import (
	"os"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

func TestMain(m *testing.M) {
	unit.InitializeAPISpecPaths()
	exitCode := m.Run()
	os.Exit(exitCode)
}

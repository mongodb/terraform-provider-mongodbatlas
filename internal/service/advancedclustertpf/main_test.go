package advancedclustertpf_test

import (
	"os"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestMain(m *testing.M) {
	// Force digest auth for unit tests so provider validation passes without OAuth2
	_ = os.Setenv("MONGODB_ATLAS_PUBLIC_API_KEY", "dummy")
	_ = os.Setenv("MONGODB_ATLAS_PRIVATE_API_KEY", "dummy")
	// Ensure Service Account auth is not selected during tests
	_ = os.Unsetenv("MONGODB_ATLAS_CLIENT_ID")
	_ = os.Unsetenv("MONGODB_ATLAS_CLIENT_SECRET")
	_ = os.Unsetenv("TF_VAR_CLIENT_ID")
	_ = os.Unsetenv("TF_VAR_CLIENT_SECRET")

	cleanup := acc.SetupSharedResources()
	exitCode := m.Run()
	cleanup()
	os.Exit(exitCode)
}

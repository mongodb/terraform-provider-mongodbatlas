package advancedclustertpf_test

import (
	"os"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestMain(m *testing.M) {
	// Only modify credentials for unit tests. Preserve GH Actions env in acceptance (TF_ACC=1).
	if os.Getenv("TF_ACC") == "" {
		// If no credentials are provided, force digest auth to satisfy provider validation.
		_ = os.Setenv("MONGODB_ATLAS_PUBLIC_API_KEY", "dummy")
		_ = os.Setenv("MONGODB_ATLAS_PRIVATE_API_KEY", "dummy")
		// Ensure Service Account auth is not selected if we just set digest
		_ = os.Unsetenv("MONGODB_ATLAS_CLIENT_ID")
		_ = os.Unsetenv("MONGODB_ATLAS_CLIENT_SECRET")
		_ = os.Unsetenv("TF_VAR_CLIENT_ID")
		_ = os.Unsetenv("TF_VAR_CLIENT_SECRET")
	}

	cleanup := acc.SetupSharedResources()
	exitCode := m.Run()
	cleanup()
	os.Exit(exitCode)
}

package mig

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

// ProjectIDGlobal returns a common global project to be used by migration tests.
// As there is a small number of mig tests this project won't hit limits.
// When `MONGODB_ATLAS_PROJECT_ID` is defined, it is used instead of creating a project. This is useful for local execution but not intended for CI executions.
func ProjectIDGlobal(tb testing.TB) string {
	tb.Helper()
	return acc.ProjectID(tb, acc.PrefixProjectKeep+"-global-mig")
}

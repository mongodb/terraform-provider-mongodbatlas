package mig

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

// ProjectIDGlobal returns a common global project to be used by migration tests.
// As there is a small number of mig tests this project won't hit limits.
func ProjectIDGlobal(tb testing.TB) string {
	tb.Helper()
	return acc.ProjectID(tb, acc.PrefixProjectKeep+"-global-mig")
}

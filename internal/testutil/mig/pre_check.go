package mig

import (
	"os"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func PreCheckBasic(tb testing.TB) {
	checkLastVersion(tb)
	acc.PreCheckBasic(tb)
}

func PreCheck(tb testing.TB) {
	checkLastVersion(tb)
	acc.PreCheck(tb)
}

func PreCheckBasicOwnerID(tb testing.TB) {
	PreCheckBasic(tb)
}

func checkLastVersion(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_LAST_VERSION") == "" {
		tb.Fatal("`MONGODB_ATLAS_LAST_VERSION` must be set for migration acceptance testing")
	}
}

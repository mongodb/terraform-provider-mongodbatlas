package mig

import (
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

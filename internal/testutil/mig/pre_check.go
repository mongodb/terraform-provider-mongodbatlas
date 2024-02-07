package mig

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func PreCheckBasic(tb testing.TB) {
	tb.Helper()
	checkLastVersion(tb)
	acc.PreCheckBasic(tb)
}

func PreCheck(tb testing.TB) {
	tb.Helper()
	checkLastVersion(tb)
	acc.PreCheck(tb)
}

func PreCheckBasicOwnerID(tb testing.TB) {
	tb.Helper()
	PreCheckBasic(tb)
}

func PreCheckCert(tb testing.TB) {
	tb.Helper()
	checkLastVersion(tb)
	acc.PreCheckCert(tb)
}

func PreCheckAtlasUsername(tb testing.TB) {
	tb.Helper()
	checkLastVersion(tb)
	acc.PreCheckAtlasUsername(tb)
}

func PreCheckPrivateEndpointServiceDataFederationOnlineArchiveRun(tb testing.TB) {
	tb.Helper()
	checkLastVersion(tb)
	acc.PreCheckPrivateEndpointServiceDataFederationOnlineArchiveRun(tb)
}

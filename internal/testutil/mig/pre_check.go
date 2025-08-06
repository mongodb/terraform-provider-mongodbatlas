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

// PreCheckBasicSleep is a helper function to call SerialSleep, see its help for more info.
// Some examples of use are when the test is calling ProjectIDExecution or GetClusterInfo to create clusters.
func PreCheckBasicSleep(tb testing.TB) func() {
	tb.Helper()
	return func() {
		PreCheckBasic(tb)
		acc.SerialSleep(tb)
	}
}

func PreCheckLast1XVersion(tb testing.TB) func() {
	tb.Helper()
	return func() {
		acc.PreCheckLast1XVersion(tb)
	}
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

func PreCheckLDAP(tb testing.TB) {
	tb.Helper()
	checkLastVersion(tb)
	acc.PreCheckLDAP(tb)
}

func PreCheckAtlasUsername(tb testing.TB) {
	tb.Helper()
	checkLastVersion(tb)
	acc.PreCheckAtlasUsername(tb)
}

func PreCheckPrivateEndpoint(tb testing.TB) {
	tb.Helper()
	checkLastVersion(tb)
	acc.PreCheckPrivateEndpoint(tb)
}

func PreCheckPeeringEnvAWS(tb testing.TB) {
	tb.Helper()
	checkLastVersion(tb)
	acc.PreCheckPeeringEnvAWS(tb)
}

func PreCheckAwsEnvBasic(tb testing.TB) {
	tb.Helper()
	checkLastVersion(tb)
	acc.PreCheckAwsEnvBasic(tb)
}

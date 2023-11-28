package acc

import (
	"os"
	"testing"
)

func DebugVersion(tb testing.TB, method, others string) {
	tb.Logf("DebugVersion METHOD: %s, others: %s", method, others)
	tb.Logf("TF_ACC_TERRAFORM_PATH: %s", os.Getenv("TF_ACC_TERRAFORM_PATH"))
	tb.Logf("TF_ACC_TERRAFORM_VERSION: %s", os.Getenv("TF_ACC_TERRAFORM_VERSION"))
}

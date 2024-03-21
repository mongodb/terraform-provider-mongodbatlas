package acc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

const (
	prefixName        = "test-acc-tf"
	prefixProject     = prefixName + "-p"
	PrefixProjectKeep = prefixProject + "-keep"
	prefixCluster     = prefixName + "-c"
	prefixIAMRole     = "mongodb-atlas-" + prefixName
	prefixIAMUser     = "arn:aws:iam::358363220050:user/mongodb-aws-iam-auth-test-user"
)

func RandomName() string {
	return acctest.RandomWithPrefix(prefixName)
}

func RandomProjectName() string {
	return acctest.RandomWithPrefix(prefixProject)
}

func RandomClusterName() string {
	return acctest.RandomWithPrefix(prefixCluster)
}

func RandomIAMRole() string {
	return acctest.RandomWithPrefix(prefixIAMRole)
}

func RandomIAMUser() string {
	return acctest.RandomWithPrefix(prefixIAMUser)
}

func RandomIP(a, b, c byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, acctest.RandIntRange(0, 255))
}

func RandomEmail() string {
	return fmt.Sprintf("%s-%s@mongodb.com", prefixName, acctest.RandString(10))
}

func RandomLDAPName() string {
	return fmt.Sprintf("CN=%s-%s@example.com,OU=users,DC=example,DC=com", prefixName, acctest.RandString(10))
}

// ProjectIDGlobal returns a common global project.
// Consider mig.ProjectIDGlobal for mig tests.
// When `MONGODB_ATLAS_PROJECT_ID` is defined, it is used instead of creating a project. This is useful for local execution but not intended for CI executions.
func ProjectIDGlobal(tb testing.TB) string {
	tb.Helper()
	return ProjectID(tb, PrefixProjectKeep+"-global")
}

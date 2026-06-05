package acc

import (
	"fmt"
	"sync/atomic"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

const (
	prefixName     = "test-acc-tf"
	prefixProject  = prefixName + "-p"
	prefixCluster  = prefixName + "-c"
	prefixStream   = prefixName + "-s"
	prefixIAMRole  = "mongodb-atlas-" + prefixName
	prefixIAMUser  = "arn:aws:iam::358363220050:user/mongodb-aws-iam-auth-test-user"
	prefixS3Bucket = "mongodb-atlas-tf"
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

func RandomStreamInstanceName() string {
	return acctest.RandomWithPrefix(prefixStream)
}

func RandomIAMRole() string {
	return acctest.RandomWithPrefix(prefixIAMRole)
}

func RandomIAMUser() string {
	return acctest.RandomWithPrefix(prefixIAMUser)
}

// ipCounter provides unique IPs per process; each run gets its own project so no cross-process seed is needed.
var ipCounter atomic.Uint32

func RandomIP() string {
	n := ipCounter.Add(1)
	return fmt.Sprintf("179.%d.%d.%d", byte(n>>16), byte(n>>8), byte(n)) //nolint:gosec // intentional octet extraction
}

func RandomEmail() string {
	return fmt.Sprintf("%s-%s@mongodb.com", prefixName, acctest.RandString(10))
}

func RandomLDAPName() string {
	return fmt.Sprintf("CN=%s-%s@example.com,OU=users,DC=example,DC=com", prefixName, acctest.RandString(10))
}

func RandomBucketName() string {
	return fmt.Sprintf("%s-%s", prefixS3Bucket, acctest.RandString(10))
}

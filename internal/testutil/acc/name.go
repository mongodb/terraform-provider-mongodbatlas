package acc

import (
	"fmt"

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

func RandomIP(a, b, c byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, acctest.RandIntRange(0, 255))
}

func RandomEmail() string {
	return fmt.Sprintf("%s-%s@mongodb.com", prefixName, acctest.RandString(10))
}

func RandomLDAPName() string {
	return fmt.Sprintf("CN=%s-%s@example.com,OU=users,DC=example,DC=com", prefixName, acctest.RandString(10))
}

func RandomS3BucketName() string {
	return fmt.Sprintf("%s-%s", prefixS3Bucket, acctest.RandString(10))
}

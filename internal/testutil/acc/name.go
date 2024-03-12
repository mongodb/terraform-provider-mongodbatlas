package acc

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

const (
	prefixName        = "test-acc-tf"
	prefixProject     = prefixName + "-p"
	prefixProjectKeep = prefixProject + "-keep"
	prefixCluster     = prefixName + "-c"
	prefixIAMRole     = "mongodb-atlas-" + prefixName
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

func RandomIP(a, b, c byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, acctest.RandIntRange(0, 255))
}

func RandomEmail() string {
	return fmt.Sprintf("%s-%s@mongodb.com", prefixName, acctest.RandString(10))
}

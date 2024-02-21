package acc

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

const (
	prefixName    = "test-acc-tf"
	prefixProject = prefixName + "-p"
	prefixCluster = prefixName + "-c"
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

func RandomIP(a, b, c byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, acctest.RandIntRange(0, 255))
}

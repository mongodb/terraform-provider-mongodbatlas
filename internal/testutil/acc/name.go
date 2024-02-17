package acc

import (
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

package acc

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/stretchr/testify/require"
)

const (
	prefixName        = "test-acc-tf"
	prefixProject     = prefixName + "-p"
	prefixProjectKeep = prefixProject + "-keep-"
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

func ProjectIDPrivateEndpoint(tb testing.TB) string {
	tb.Helper()
	return projectID(tb, prefixProjectKeep+"private-endpoint")
}

func projectID(tb testing.TB, name string) string {
	tb.Helper()
	resp, _, _ := ConnV2().ProjectsApi.GetProjectByName(context.Background(), name).Execute()
	id := resp.GetId()
	require.NotEmpty(tb, id)
	return id
}

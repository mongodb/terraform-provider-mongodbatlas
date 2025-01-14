package acc_test

import (
	"fmt"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
)

func Test_NextProjectIDClusterName(t *testing.T) {
	projectIDReturner := func(name string) string {
		return fmt.Sprintf("%s-name", name)
	}
	projectIDs := map[string]int{}
	clusterNames := map[string]int{}

	addProjectIDClusterName := func() {
		projectID, clusterName := acc.NextProjectIDClusterName(projectIDReturner)
		projectIDs[projectID]++
		clusterNames[clusterName]++
	}
	for range acc.MaxClustersPerProject {
		addProjectIDClusterName()
	}
	assert.Len(t, projectIDs, 1)
	assert.Len(t, clusterNames, acc.MaxClustersPerProject)
	addProjectIDClusterName()
	assert.Len(t, projectIDs, 2)
	assert.Len(t, clusterNames, acc.MaxClustersPerProject+1)
}

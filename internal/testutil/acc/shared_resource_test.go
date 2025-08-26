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

	addProjectIDClusterName := func(nodeCount int, freeTierClusterCount int) {
		projectID, clusterName := acc.NextProjectIDClusterName(nodeCount, freeTierClusterCount, projectIDReturner)
		projectIDs[projectID]++
		clusterNames[clusterName]++
	}
	for range acc.MaxClusterNodesPerProject {
		addProjectIDClusterName(1, 0)
	}
	assert.Len(t, projectIDs, 1)
	assert.Len(t, clusterNames, acc.MaxClusterNodesPerProject)
	addProjectIDClusterName(1, 0)
	assert.Len(t, projectIDs, 2)
	assert.Len(t, clusterNames, acc.MaxClusterNodesPerProject+1)
	addProjectIDClusterName(acc.MaxClusterNodesPerProject, 0)
	assert.Len(t, projectIDs, 3)
	assert.Len(t, clusterNames, acc.MaxClusterNodesPerProject+2)
	addProjectIDClusterName(1, 0)
	assert.Len(t, projectIDs, 4)
	assert.Len(t, clusterNames, acc.MaxClusterNodesPerProject+3)
	addProjectIDClusterName(0, 1) // adds free tier, shares existing project
	assert.Len(t, projectIDs, 4)
	assert.Len(t, clusterNames, acc.MaxClusterNodesPerProject+4)
	addProjectIDClusterName(0, 1) // second free tier cluster creates a new project
	assert.Len(t, projectIDs, 5)
	assert.Len(t, clusterNames, acc.MaxClusterNodesPerProject+5)
}

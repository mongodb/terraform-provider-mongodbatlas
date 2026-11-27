package networkcontainer_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312028/admin"
	"go.mongodb.org/atlas-sdk/v20250312028/mockadmin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/networkcontainer"
)

const (
	testProjectID   = "project-id"
	testContainerID = "container-id"
)

func containerResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]any{})
	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   testProjectID,
		"container_id": testContainerID,
	}))
	return d
}

func TestResourceRefreshFunc_pendingReturnsNonNilResult(t *testing.T) {
	// StateChangeConf counts a nil result as a not-found check and aborts after 21 of them,
	// ignoring its timeout, so pending polls must return a non-nil result.
	ctx := context.Background()
	client := func(containerProvisioned bool, deleteErr error) *admin.APIClient {
		api := mockadmin.NewNetworkPeeringAPI(t)
		container := &admin.CloudProviderContainer{Id: new(testContainerID), Provisioned: new(containerProvisioned)}
		api.EXPECT().GetGroupContainer(mock.Anything, testProjectID, testContainerID).
			Return(admin.GetGroupContainerApiRequest{ApiService: api}).Once()
		api.EXPECT().GetGroupContainerExecute(mock.Anything).Return(container, nil, nil).Once()
		if !containerProvisioned {
			api.EXPECT().DeleteGroupContainer(mock.Anything, testProjectID, testContainerID).
				Return(admin.DeleteGroupContainerApiRequest{ApiService: api}).Once()
			api.EXPECT().DeleteGroupContainerExecute(mock.Anything).Return(nil, deleteErr).Once()
		}
		return &admin.APIClient{NetworkPeeringAPI: api}
	}

	t.Run("provisioned container", func(t *testing.T) {
		result, state, err := networkcontainer.NewRefreshFuncForTest(ctx, containerResourceData(t), client(true, nil))()
		require.NoError(t, err)
		assert.Equal(t, "provisioned_container", state)
		assert.NotNil(t, result)
	})

	t.Run("delete blocked by a peering in use", func(t *testing.T) {
		client := client(false, errors.New(`HTTP 409 Conflict (Error code: "CONTAINERS_IN_USE")`))
		result, state, err := networkcontainer.NewRefreshFuncForTest(ctx, containerResourceData(t), client)()
		require.NoError(t, err)
		assert.Equal(t, "provisioned_container", state)
		assert.NotNil(t, result)
	})
}

func TestResourceRefreshFunc_containerDeleted(t *testing.T) {
	api := mockadmin.NewNetworkPeeringAPI(t)
	api.EXPECT().GetGroupContainer(mock.Anything, testProjectID, testContainerID).
		Return(admin.GetGroupContainerApiRequest{ApiService: api}).Once()
	api.EXPECT().GetGroupContainerExecute(mock.Anything).
		Return(nil, &http.Response{StatusCode: http.StatusNotFound}, errors.New("not found")).Once()

	result, state, err := networkcontainer.NewRefreshFuncForTest(
		context.Background(), containerResourceData(t), &admin.APIClient{NetworkPeeringAPI: api},
	)()
	require.NoError(t, err)
	assert.Equal(t, "deleted", state)
	assert.NotNil(t, result)
}

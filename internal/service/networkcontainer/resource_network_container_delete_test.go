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

func apiWithContainerGet(t *testing.T, getStatus int, getErr error) *mockadmin.NetworkPeeringAPI {
	t.Helper()
	api := mockadmin.NewNetworkPeeringAPI(t)
	api.EXPECT().GetGroupContainer(mock.Anything, testProjectID, testContainerID).
		Return(admin.GetGroupContainerApiRequest{ApiService: api}).Once()
	api.EXPECT().GetGroupContainerExecute(mock.Anything).
		Return(&admin.CloudProviderContainer{Id: new(testContainerID)}, &http.Response{StatusCode: getStatus}, getErr).Once()
	return api
}

func TestResourceDeleteRefreshFunc(t *testing.T) {
	ctx := context.Background()

	t.Run("container already gone before delete", func(t *testing.T) {
		api := apiWithContainerGet(t, http.StatusNotFound, errors.New("not found"))
		result, state, err := networkcontainer.NewDeleteRefreshFuncForTest(ctx, containerResourceData(t), &admin.APIClient{NetworkPeeringAPI: api})()
		require.NoError(t, err)
		assert.Equal(t, "deleted", state)
		assert.NotNil(t, result)
	})

	t.Run("container in use by a peering", func(t *testing.T) {
		api := apiWithContainerGet(t, http.StatusOK, nil)
		api.EXPECT().DeleteGroupContainer(mock.Anything, testProjectID, testContainerID).
			Return(admin.DeleteGroupContainerApiRequest{ApiService: api}).Once()
		api.EXPECT().DeleteGroupContainerExecute(mock.Anything).
			Return(&http.Response{StatusCode: http.StatusConflict}, errors.New(`HTTP 409 Conflict (Error code: "CONTAINERS_IN_USE")`)).Once()

		result, state, err := networkcontainer.NewDeleteRefreshFuncForTest(ctx, containerResourceData(t), &admin.APIClient{NetworkPeeringAPI: api})()
		require.NoError(t, err)
		assert.Equal(t, "provisioned_container", state)
		assert.NotNil(t, result)
	})

	t.Run("delete fails with a non-retryable error", func(t *testing.T) {
		api := apiWithContainerGet(t, http.StatusOK, nil)
		api.EXPECT().DeleteGroupContainer(mock.Anything, testProjectID, testContainerID).
			Return(admin.DeleteGroupContainerApiRequest{ApiService: api}).Once()
		api.EXPECT().DeleteGroupContainerExecute(mock.Anything).
			Return(&http.Response{StatusCode: http.StatusForbidden}, errors.New("forbidden")).Once()

		_, _, err := networkcontainer.NewDeleteRefreshFuncForTest(ctx, containerResourceData(t), &admin.APIClient{NetworkPeeringAPI: api})()
		require.ErrorContains(t, err, "forbidden")
	})

	t.Run("container deleted", func(t *testing.T) {
		api := apiWithContainerGet(t, http.StatusOK, nil)
		api.EXPECT().DeleteGroupContainer(mock.Anything, testProjectID, testContainerID).
			Return(admin.DeleteGroupContainerApiRequest{ApiService: api}).Once()
		api.EXPECT().DeleteGroupContainerExecute(mock.Anything).
			Return(&http.Response{StatusCode: http.StatusNoContent}, nil).Once()

		_, state, err := networkcontainer.NewDeleteRefreshFuncForTest(ctx, containerResourceData(t), &admin.APIClient{NetworkPeeringAPI: api})()
		require.NoError(t, err)
		assert.Equal(t, "deleted", state)
	})
}

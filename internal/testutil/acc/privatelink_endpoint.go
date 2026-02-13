package acc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privatelinkendpoint"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
)

func createPrivateLinkEndpoint(tb testing.TB, projectID, providerName, region string) string {
	tb.Helper()
	request := &admin.CloudProviderEndpointServiceRequest{
		ProviderName: providerName,
		Region:       region,
	}
	privateEndpoint, _, err := ConnV2().PrivateEndpointServicesApi.CreatePrivateEndpointService(tb.Context(), projectID, request).Execute()
	require.NoError(tb, err)
	stateConf := privatelinkendpoint.CreateStateChangeConfig(tb.Context(), ConnV2(), projectID, providerName, privateEndpoint.GetId(), 1*time.Hour)
	_, err = stateConf.WaitForStateContext(tb.Context())
	require.NoError(tb, err, "Private link endpoint creation failed: %s, err: %s", privateEndpoint.GetId(), err)
	return privateEndpoint.GetId()
}

func deletePrivateLinkEndpoint(projectID, providerName, privateLinkEndpointID string) error {
	const maxConflictRetries = 3
	for i := range maxConflictRetries {
		resp, err := ConnV2().PrivateEndpointServicesApi.DeletePrivateEndpointService(context.Background(), projectID, providerName, privateLinkEndpointID).Execute()
		if err != nil {
			// 409 Conflict occurs when an attached endpoint is still being removed.
			// Always happens for deleteOnCreateTimeout tests since the resource does not wait for deletion in that case.
			if validate.StatusConflict(resp) && i < maxConflictRetries-1 {
				fmt.Printf("Private link endpoint deletion failed, will retry in 10s: %s, error: %s\n", privateLinkEndpointID, err)
				time.Sleep(10 * time.Second)
				continue
			}
			return fmt.Errorf("failed to delete private link endpoint %s: %w", privateLinkEndpointID, err)
		}
		break
	}
	stateConf := privatelinkendpoint.DeleteStateChangeConfig(context.Background(), ConnV2(), projectID, providerName, privateLinkEndpointID, 1*time.Hour)
	_, err := stateConf.WaitForStateContext(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete private link endpoint %s: %w", privateLinkEndpointID, err)
	}
	return nil
}

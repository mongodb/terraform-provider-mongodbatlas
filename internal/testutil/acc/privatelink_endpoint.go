package acc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privatelinkendpoint"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
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

func deletePrivateLinkEndpoint(projectID, providerName, privateLinkEndpointID string) {
	_, err := ConnV2().PrivateEndpointServicesApi.DeletePrivateEndpointService(context.Background(), projectID, providerName, privateLinkEndpointID).Execute()
	if err != nil {
		fmt.Printf("Failed to delete private link endpoint %s: %s\n", privateLinkEndpointID, err)
		return
	}
	stateConf := privatelinkendpoint.DeleteStateChangeConfig(context.Background(), ConnV2(), projectID, providerName, privateLinkEndpointID, 1*time.Hour)
	_, err = stateConf.WaitForStateContext(context.Background())
	if err != nil {
		fmt.Printf("Failed to delete private link endpoint %s: %s\n", privateLinkEndpointID, err)
	}
}

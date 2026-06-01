package acc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/clientcredentials"
)

const azureManagementAPIVersion = "2024-01-01"

// CleanAzurePeeringConnections deletes any dangling Azure VNet peering connections
// left by previous test runs. Called in PreCheck of Azure network peering tests.
// Silently returns if required env vars are not set.
func CleanAzurePeeringConnections(tb testing.TB) {
	tb.Helper()
	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_APP_SECRET")
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	resourceGroupName := os.Getenv("AZURE_RESOURCE_GROUP_NAME")
	if tenantID == "" || clientID == "" || clientSecret == "" || subscriptionID == "" || resourceGroupName == "" {
		return
	}
	conf := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID),
		Scopes:       []string{"https://management.azure.com/.default"},
	}
	httpClient := conf.Client(context.Background())
	for _, vnetName := range []string{os.Getenv("AZURE_VNET_NAME"), os.Getenv("AZURE_VNET_NAME_UPDATED")} {
		if vnetName == "" {
			continue
		}
		cleanPeeringsForVNet(tb, httpClient, subscriptionID, resourceGroupName, vnetName)
	}
}

func cleanPeeringsForVNet(tb testing.TB, client *http.Client, subscriptionID, resourceGroupName, vnetName string) {
	tb.Helper()
	baseURL := fmt.Sprintf(
		"https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/virtualNetworkPeerings",
		subscriptionID, resourceGroupName, vnetName,
	)
	listReq, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseURL+"?api-version="+azureManagementAPIVersion, http.NoBody)
	require.NoError(tb, err)
	listResp, err := client.Do(listReq)
	require.NoError(tb, err)
	defer listResp.Body.Close()
	require.Equal(tb, http.StatusOK, listResp.StatusCode, "list peerings for vnet %s", vnetName)
	var listResult struct {
		Value []struct {
			Name string `json:"name"`
		} `json:"value"`
	}
	if err := json.NewDecoder(listResp.Body).Decode(&listResult); err != nil || len(listResult.Value) == 0 {
		return
	}
	for _, p := range listResult.Value {
		deleteURL := fmt.Sprintf("%s/%s?api-version=%s", baseURL, p.Name, azureManagementAPIVersion)
		delReq, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, deleteURL, http.NoBody)
		require.NoError(tb, err)
		delResp, err := client.Do(delReq)
		require.NoError(tb, err)
		delResp.Body.Close()
		require.Less(tb, delResp.StatusCode, 300, "delete peering %s in vnet %s", p.Name, vnetName)
	}
}

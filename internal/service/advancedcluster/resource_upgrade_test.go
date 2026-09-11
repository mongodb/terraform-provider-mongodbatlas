package advancedcluster_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
	"go.mongodb.org/atlas-sdk/v20250312024/mockadmin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
)

func TestAdvancedClusterUpdate_FlexUpgradeDatabaseEdition(t *testing.T) {
	testCases := map[string]struct {
		stateEdition *string
		planEdition  *string
	}{
		"unchanged CORE":     {stateEdition: new("CORE"), planEdition: new("CORE")},
		"unchanged INFINITE": {stateEdition: new("INFINITE"), planEdition: new("INFINITE")},
		"omitted":            {},
		"newly configured":   {planEdition: new("CORE")},
		"changed":            {stateEdition: new("CORE"), planEdition: new("INFINITE")},
		"removed":            {stateEdition: new("CORE")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			const projectID = "111111111111111111111111"
			api := mockadmin.NewFlexClustersAPI(t)
			// Stop at the request boundary before polling for an Atlas state transition.
			upgradeErr := errors.New("stop after upgrade request")
			api.EXPECT().TenantUpgrade(mock.Anything, projectID, mock.Anything).
				Run(func(_ context.Context, _ string, req *admin.AtlasTenantClusterUpgradeRequest20240805) {
					payload, err := json.Marshal(req)
					require.NoError(t, err)
					var fields map[string]any
					require.NoError(t, json.Unmarshal(payload, &fields))
					if tc.planEdition == nil {
						assert.NotContains(t, fields, "databaseEdition")
					} else {
						assert.Equal(t, *tc.planEdition, fields["databaseEdition"])
					}
				}).Return(admin.TenantUpgradeApiRequest{ApiService: api}).Once()
			api.EXPECT().TenantUpgradeExecute(mock.Anything).Return(nil, nil, upgradeErr).Once()

			rs, resourceSchema := configuredClusterResource(t, &admin.APIClient{FlexClustersAPI: api})
			state := clusterUpgradeState(t, resourceSchema, projectID, "FLEX", tc.stateEdition)
			plan := clusterUpgradeState(t, resourceSchema, projectID, "AWS", tc.planEdition)
			resp := &resource.UpdateResponse{State: state}
			rs.Update(t.Context(), resource.UpdateRequest{
				State: state,
				Plan:  tfsdk.Plan(plan),
			}, resp)
			require.Len(t, resp.Diagnostics, 1)
			assert.Equal(t, "Error in flex upgrade", resp.Diagnostics[0].Summary())
			assert.Contains(t, resp.Diagnostics[0].Detail(), upgradeErr.Error())
		})
	}
}

func TestAdvancedClusterUpdate_DatabaseEditionRequests(t *testing.T) {
	testCases := map[string]struct {
		stateEdition *string
		planEdition  *string
		dedicated    bool
		twoNodes     bool
	}{
		"free upgrade omitted":            {},
		"free upgrade CORE":               {planEdition: new("CORE")},
		"free upgrade INFINITE":           {planEdition: new("INFINITE")},
		"free upgrade two-node INFINITE":  {planEdition: new("INFINITE"), twoNodes: true},
		"free upgrade unchanged CORE":     {stateEdition: new("CORE"), planEdition: new("CORE")},
		"free upgrade unchanged INFINITE": {stateEdition: new("INFINITE"), planEdition: new("INFINITE")},
		"free upgrade changed":            {stateEdition: new("CORE"), planEdition: new("INFINITE")},
		"free upgrade removed":            {stateEdition: new("CORE")},
		"dedicated update changed":        {stateEdition: new("CORE"), planEdition: new("INFINITE"), dedicated: true},
		"dedicated update removed":        {stateEdition: new("CORE"), dedicated: true},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			const projectID = "111111111111111111111111"
			responseBody := &clusterRequestBody{Reader: strings.NewReader(`{"error":400,"errorCode":"DATABASE_EDITION_TEST_ERROR","detail":"databaseEdition cannot be changed","reason":"Bad Request"}`)}
			calls := 0
			client, err := admin.NewClient(admin.UseBaseURL("https://example.invalid"), admin.UseHTTPClient(&http.Client{
				Transport: clusterRequestTransport(func(req *http.Request) (*http.Response, error) {
					defer req.Body.Close()
					calls++
					expectedMethod, expectedPath, expectedVersion := http.MethodPost, "/api/atlas/v2/groups/"+projectID+"/clusters/tenantUpgrade", "application/vnd.atlas.2023-01-01+json"
					if tc.dedicated {
						expectedMethod, expectedPath, expectedVersion = http.MethodPatch, "/api/atlas/v2/groups/"+projectID+"/clusters/test", "application/vnd.atlas.2024-10-23+json"
					}
					assert.Equal(t, expectedMethod, req.Method)
					assert.Equal(t, expectedPath, req.URL.Path)
					assert.Equal(t, expectedVersion, req.Header.Get("Accept"))
					assert.Equal(t, expectedVersion, req.Header.Get("Content-Type"))
					var payload map[string]any
					require.NoError(t, json.NewDecoder(req.Body).Decode(&payload))
					switch {
					case tc.planEdition != nil:
						assert.Equal(t, *tc.planEdition, payload["databaseEdition"])
					case tc.dedicated:
						assert.Contains(t, payload, "databaseEdition")
						assert.Nil(t, payload["databaseEdition"])
					default:
						assert.NotContains(t, payload, "databaseEdition")
					}
					if !tc.dedicated {
						assert.Equal(t, "test", payload["name"])
						assert.Equal(t, true, payload["providerBackupEnabled"])
						assert.Equal(t, map[string]any{"providerName": "AWS", "regionName": "US_EAST_1", "instanceSizeName": "M10"}, payload["providerSettings"])
						assert.NotContains(t, payload, "backupEnabled")
						assert.NotContains(t, payload, "replicationSpecs")
						if tc.twoNodes {
							assert.InDelta(t, 2, payload["replicationFactor"], 0)
						} else {
							assert.NotContains(t, payload, "replicationFactor")
						}
					}
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       responseBody,
						Request:    req,
					}, nil
				}),
			}))
			require.NoError(t, err)
			rs, resourceSchema := configuredClusterResource(t, client)
			sourceProvider, expectedSummary := "TENANT", "Error in tenant upgrade"
			if tc.dedicated {
				sourceProvider, expectedSummary = "AWS", "Error in update"
			}
			state := clusterUpgradeState(t, resourceSchema, projectID, sourceProvider, tc.stateEdition)
			plan := clusterUpgradeState(t, resourceSchema, projectID, "AWS", tc.planEdition)
			if tc.twoNodes {
				diags := plan.SetAttribute(t.Context(), path.Root("replication_specs").AtListIndex(0).AtName("region_configs").AtListIndex(0).AtName("electable_specs").AtName("node_count"), 2)
				require.False(t, diags.HasError(), diags)
			}
			resp := &resource.UpdateResponse{State: state}
			rs.Update(t.Context(), resource.UpdateRequest{
				State: state,
				Plan:  tfsdk.Plan(plan),
			}, resp)
			require.Equal(t, 1, calls)
			require.Len(t, resp.Diagnostics, 1)
			assert.Equal(t, expectedSummary, resp.Diagnostics[0].Summary())
			assert.Contains(t, resp.Diagnostics[0].Detail(), "DATABASE_EDITION_TEST_ERROR")
			assert.True(t, responseBody.closed)
		})
	}
}

type clusterRequestTransport func(*http.Request) (*http.Response, error)

func (transport clusterRequestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return transport(req)
}

type clusterRequestBody struct {
	io.Reader
	closed bool
}

func (body *clusterRequestBody) Close() error {
	body.closed = true
	return nil
}

func configuredClusterResource(t *testing.T, client *admin.APIClient) (resource.Resource, schema.Schema) {
	t.Helper()
	rs := config.AnalyticsResourceFunc(advancedcluster.Resource())()
	rsConfigure, ok := rs.(resource.ResourceWithConfigure)
	require.True(t, ok)
	configureResp := &resource.ConfigureResponse{}
	rsConfigure.Configure(t.Context(), resource.ConfigureRequest{
		ProviderData: &config.MongoDBClient{AtlasV2: client},
	}, configureResp)
	require.False(t, configureResp.Diagnostics.HasError(), configureResp.Diagnostics)
	schemaResp := &resource.SchemaResponse{}
	rs.Schema(t.Context(), resource.SchemaRequest{}, schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError(), schemaResp.Diagnostics)
	return rs, schemaResp.Schema
}

func clusterUpgradeState(t *testing.T, resourceSchema schema.Schema, projectID, providerName string, databaseEdition *string) tfsdk.State {
	t.Helper()
	region := map[string]any{
		"provider_name": providerName,
		"region_name":   "US_EAST_1",
		"priority":      7,
	}
	switch providerName {
	case "FLEX":
		region["backing_provider_name"] = "AWS"
	case "TENANT":
		region["backing_provider_name"] = "AWS"
		region["electable_specs"] = map[string]any{"instance_size": "M0"}
	default:
		region["electable_specs"] = map[string]any{"instance_size": "M10", "node_count": 3}
	}
	payload, err := json.Marshal(map[string]any{
		"project_id":       projectID,
		"name":             "test",
		"cluster_type":     "REPLICASET",
		"backup_enabled":   providerName == "AWS",
		"database_edition": databaseEdition,
		"replication_specs": []any{map[string]any{
			"region_configs": []any{region},
		}},
	})
	require.NoError(t, err)
	raw, err := (tfprotov6.DynamicValue{JSON: payload}).Unmarshal(resourceSchema.Type().TerraformType(t.Context()))
	require.NoError(t, err)
	return tfsdk.State{Schema: resourceSchema, Raw: raw}
}

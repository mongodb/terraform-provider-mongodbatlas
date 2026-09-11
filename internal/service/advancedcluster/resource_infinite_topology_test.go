package advancedcluster_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312024/admin"
	"go.mongodb.org/atlas-sdk/v20250312024/mockadmin"

	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

func TestInfiniteClusterTypeValidation(t *testing.T) {
	for _, name := range []string{"ASSUME_ROLE_ARN", "TF_VAR_ASSUME_ROLE_ARN", "MONGODB_ATLAS_CLIENT_ID", "MONGODB_ATLAS_CLIENT_SECRET"} {
		t.Setenv(name, "")
	}
	for _, tc := range []struct {
		clusterType, edition string
		unsupported          bool
	}{
		{"SHARDED", "INFINITE", true},
		{"GEOSHARDED", "INFINITE", true},
		{"REPLICASET", "INFINITE", false},
		{"SHARDED", "CORE", false},
		{"GEOSHARDED", "CORE", false},
		{"SHARDED", "", false},
	} {
		t.Run(tc.clusterType+"/"+tc.edition, func(t *testing.T) {
			var expected *regexp.Regexp
			if tc.unsupported {
				expected = regexp.MustCompile(`(?s)Unsupported INFINITE cluster type.*cluster_type = "` + tc.clusterType + `".*not supported.*newer provider version`)
			}
			transport := &noTopologyRequests{}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, transport),
				Steps: []resource.TestStep{{
					Config: topologyValidationConfig(tc.clusterType, tc.edition), PlanOnly: true,
					ExpectError: expected, ExpectNonEmptyPlan: !tc.unsupported,
				}},
			})
			require.Zero(t, transport.calls.Load(), "validation and planning must not call Atlas")
		})
	}
}

func TestInfiniteClusterTypeUnknownPlanAndApply(t *testing.T) {
	ctx := t.Context()
	for _, clusterType := range []string{"SHARDED", "GEOSHARDED"} {
		for _, unknownAttribute := range []string{"cluster_type", "database_edition"} {
			for _, operation := range []string{"create", "update"} {
				t.Run(clusterType+"/"+unknownAttribute+"/"+operation, func(t *testing.T) {
					r := advancedcluster.Resource()
					var schemaResponse frameworkresource.SchemaResponse
					r.Schema(ctx, frameworkresource.SchemaRequest{}, &schemaResponse)
					typ := schemaResponse.Schema.Type().TerraformType(ctx)
					resolved := topologyValidationAttributes(clusterType, "INFINITE")
					unknown := topologyValidationAttributes(clusterType, "INFINITE")
					unknown[unknownAttribute] = tftypes.UnknownValue
					prior := topologyTestValue(typ, topologyValidationAttributes("REPLICASET", "CORE"))
					if operation == "create" {
						prior = tftypes.NewValue(typ, nil)
					}
					dynamic := func(value tftypes.Value) *tfprotov6.DynamicValue {
						result, err := tfprotov6.NewDynamicValue(typ, value)
						require.NoError(t, err)
						return &result
					}
					server, err := acc.TestAccProviderV6Factories["mongodbatlas"]()
					require.NoError(t, err)
					plan := func(attributes map[string]any) *tfprotov6.PlanResourceChangeResponse {
						value := dynamic(topologyTestValue(typ, attributes))
						response, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
							TypeName: "mongodbatlas_advanced_cluster", PriorState: dynamic(prior), ProposedNewState: value, Config: value,
						})
						require.NoError(t, err)
						return response
					}
					initial := plan(unknown)
					require.Empty(t, initial.Diagnostics)
					initialValue, err := initial.PlannedState.Unmarshal(typ)
					require.NoError(t, err)
					attribute, _, err := tftypes.WalkAttributePath(initialValue, tftypes.NewAttributePath().WithAttributeName(unknownAttribute))
					require.NoError(t, err)
					require.False(t, attribute.(tftypes.Value).IsKnown())
					final := plan(resolved)
					require.Len(t, final.Diagnostics, 1)
					require.Equal(t, "Unsupported INFINITE cluster type", final.Diagnostics[0].Summary)
					require.Contains(t, final.Diagnostics[0].Detail, clusterType)

					api := mockadmin.NewClustersAPI(t)
					r.(config.ImplementedResource).SetClient(&config.MongoDBClient{AtlasV2: &admin.APIClient{ClustersAPI: api}})
					resolvedPlan := tfsdk.Plan{Schema: schemaResponse.Schema, Raw: topologyTestValue(typ, resolved)}
					if operation == "create" {
						var response frameworkresource.CreateResponse
						r.Create(ctx, frameworkresource.CreateRequest{Plan: resolvedPlan}, &response)
						require.Len(t, response.Diagnostics.Errors(), 1)
						require.Equal(t, "Unsupported INFINITE cluster type", response.Diagnostics.Errors()[0].Summary())
					} else {
						var response frameworkresource.UpdateResponse
						r.Update(ctx, frameworkresource.UpdateRequest{
							Plan: resolvedPlan, State: tfsdk.State{Schema: schemaResponse.Schema, Raw: prior},
						}, &response)
						require.Len(t, response.Diagnostics.Errors(), 1)
						require.Equal(t, "Unsupported INFINITE cluster type", response.Diagnostics.Errors()[0].Summary())
					}
					require.Empty(t, api.Calls, "unsupported topologies must fail before any Atlas request")
				})
			}
		}
	}
}

func TestUpdateVerifiesImportedDatabaseEditionBeforeWrites(t *testing.T) {
	for _, clusterType := range []string{"SHARDED", "GEOSHARDED"} {
		for _, tc := range []struct {
			requestedEdition       any
			lookupError            error
			name, effectiveEdition string
		}{
			{name: "omitted edition", effectiveEdition: "INFINITE"},
			{name: "explicit CORE cannot bypass effective edition", requestedEdition: "CORE", effectiveEdition: "INFINITE"},
			{name: "CORE permits update", effectiveEdition: "CORE"},
			{name: "failed lookup prevents update", lookupError: errors.New("edition lookup failed")},
		} {
			t.Run(clusterType+"/"+tc.name, func(t *testing.T) {
				ctx := t.Context()
				r := advancedcluster.Resource()
				var schemaResponse frameworkresource.SchemaResponse
				r.Schema(ctx, frameworkresource.SchemaRequest{}, &schemaResponse)
				typ := schemaResponse.Schema.Type().TerraformType(ctx)
				prior := topologyTestValue(typ, topologyValidationAttributes("REPLICASET", nil))
				attributes := topologyValidationAttributes(clusterType, tc.requestedEdition)
				if tc.effectiveEdition == "INFINITE" {
					attributes["pinned_fcv"] = map[string]any{"expiration_date": "2099-01-01T00:00:00Z", "version": "8.0"}
				}
				plan := tfsdk.Plan{Schema: schemaResponse.Schema, Raw: topologyTestValue(typ, attributes)}
				api := mockadmin.NewClustersAPI(t)
				r.(config.ImplementedResource).SetClient(&config.MongoDBClient{AtlasV2: &admin.APIClient{ClustersAPI: api}})
				api.EXPECT().GetCluster(mock.Anything, dummyProjectID, "example").Return(admin.GetClusterApiRequest{ApiService: api}).Once()
				api.EXPECT().GetClusterExecute(mock.Anything).Return(&admin.ClusterDescription20240805{
					ClusterType: new("REPLICASET"), EffectiveDatabaseEdition: new(tc.effectiveEdition),
				}, nil, tc.lookupError).Once()
				apiError := errors.New("update request inspected")
				if tc.effectiveEdition == "CORE" {
					api.EXPECT().UpdateCluster(mock.Anything, dummyProjectID, "example", mock.Anything).Return(admin.UpdateClusterApiRequest{ApiService: api}).Once()
					api.EXPECT().UpdateClusterExecute(mock.Anything).Return(nil, nil, apiError).Once()
				}
				var response frameworkresource.UpdateResponse
				r.Update(ctx, frameworkresource.UpdateRequest{
					Plan: plan, State: tfsdk.State{Schema: schemaResponse.Schema, Raw: prior},
				}, &response)
				require.Len(t, response.Diagnostics.Errors(), 1)
				switch {
				case tc.lookupError != nil:
					require.Contains(t, response.Diagnostics.Errors()[0].Detail(), tc.lookupError.Error())
				case tc.effectiveEdition == "CORE":
					require.Contains(t, response.Diagnostics.Errors()[0].Detail(), apiError.Error())
				default:
					require.Equal(t, "Unsupported INFINITE cluster type", response.Diagnostics.Errors()[0].Summary())
					require.Contains(t, response.Diagnostics.Errors()[0].Detail(), clusterType)
				}
				if tc.effectiveEdition != "CORE" {
					api.AssertNotCalled(t, "UpdateCluster", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				}
				api.AssertNotCalled(t, "PinFeatureCompatibilityVersion", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
		}
	}
}

func TestCreateSkipsInfiniteValidationWhenEditionOmitted(t *testing.T) {
	for _, clusterType := range []string{"SHARDED", "GEOSHARDED"} {
		t.Run(clusterType, func(t *testing.T) {
			ctx := t.Context()
			r := advancedcluster.Resource()
			var schemaResponse frameworkresource.SchemaResponse
			r.Schema(ctx, frameworkresource.SchemaRequest{}, &schemaResponse)
			typ := schemaResponse.Schema.Type().TerraformType(ctx)
			plan := tfsdk.Plan{Schema: schemaResponse.Schema, Raw: topologyTestValue(typ, topologyValidationAttributes(clusterType, nil))}
			api := mockadmin.NewClustersAPI(t)
			r.(config.ImplementedResource).SetClient(&config.MongoDBClient{AtlasV2: &admin.APIClient{ClustersAPI: api}})
			apiError := errors.New("create request inspected")
			api.EXPECT().CreateCluster(mock.Anything, dummyProjectID, mock.Anything).Return(admin.CreateClusterApiRequest{ApiService: api}).Once()
			api.EXPECT().CreateClusterExecute(mock.Anything).Return(nil, nil, apiError).Once()
			var response frameworkresource.CreateResponse
			r.Create(ctx, frameworkresource.CreateRequest{Plan: plan}, &response)
			require.Len(t, response.Diagnostics.Errors(), 1)
			require.Contains(t, response.Diagnostics.Errors()[0].Detail(), apiError.Error(),
				"an omitted database_edition defaults to CORE server-side and must not be treated as INFINITE")
		})
	}
}

func TestCreateRejectsUnexpectedInfiniteTopologyPreservesState(t *testing.T) {
	t.Setenv("ASSUME_ROLE_ARN", "")
	t.Setenv("TF_VAR_ASSUME_ROLE_ARN", "")
	t.Setenv("MONGODB_ATLAS_CLIENT_ID", "")
	t.Setenv("MONGODB_ATLAS_CLIENT_SECRET", "")
	require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach())
	for _, clusterType := range []string{"SHARDED", "GEOSHARDED"} {
		t.Run(clusterType, func(t *testing.T) {
			ctx := t.Context()
			httpMock := &unexpectedInfiniteCreateHTTPMock{clusterType: clusterType}
			server, err := unit.TestAccProviderV6FactoriesWithMock(t, httpMock)["mongodbatlas"]()
			require.NoError(t, err)
			schemas, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
			require.NoError(t, err)
			require.Empty(t, schemas.Diagnostics)
			dynamic := func(typ tftypes.Type, input any) *tfprotov6.DynamicValue {
				value, err := tfprotov6.NewDynamicValue(typ, topologyTestValue(typ, input))
				require.NoError(t, err)
				return &value
			}
			providerConfig := dynamic(schemas.Provider.ValueType(), map[string]any{
				"public_key": "test-public-key", "private_key": "test-private-key", "base_url": "https://atlas.invalid/",
			})
			configured, err := server.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{Config: providerConfig})
			require.NoError(t, err)
			require.Empty(t, configured.Diagnostics)
			typ := schemas.ResourceSchemas["mongodbatlas_advanced_cluster"].ValueType()
			clusterConfig := dynamic(typ, map[string]any{
				"name": "example", "project_id": "111111111111111111111111", "cluster_type": clusterType,
				"paused": true,
				"replication_specs": []any{map[string]any{"region_configs": []any{map[string]any{
					"provider_name": "AWS", "region_name": "US_EAST_1", "priority": int64(7),
					"electable_specs": map[string]any{"instance_size": "M10", "node_count": int64(3)},
				}}}},
				"advanced_configuration": map[string]any{"minimum_enabled_tls_protocol": "TLS1_2"},
			})
			empty := dynamic(typ, nil)
			plan, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
				TypeName: "mongodbatlas_advanced_cluster", PriorState: empty, ProposedNewState: clusterConfig, Config: clusterConfig,
			})
			require.NoError(t, err)
			require.Empty(t, plan.Diagnostics)
			created, err := server.ApplyResourceChange(ctx, &tfprotov6.ApplyResourceChangeRequest{
				TypeName: "mongodbatlas_advanced_cluster", PriorState: empty, PlannedState: plan.PlannedState, Config: clusterConfig,
			})
			require.NoError(t, err)
			require.Len(t, created.Diagnostics, 1, "the unsupported-topology diagnostic must not be accompanied by invalid state errors")
			require.Equal(t, "Unsupported INFINITE cluster type", created.Diagnostics[0].Summary)
			require.Contains(t, created.Diagnostics[0].Detail, clusterType)
			require.Equal(t, []string{"POST", "GET"}, httpMock.methods, "no pause or advanced configuration writes may follow the unexpected edition")
			state, err := created.NewState.Unmarshal(typ)
			require.NoError(t, err)
			require.False(t, state.IsNull(), "a successfully created cluster must remain tracked even when its topology is rejected")
			require.True(t, state.IsFullyKnown(), "retained state must not include unknown plan values")
			for name, expected := range map[string]string{
				"name": "example", "project_id": "111111111111111111111111", "cluster_id": "333333333333333333333333",
			} {
				value, _, err := tftypes.WalkAttributePath(state, tftypes.NewAttributePath().WithAttributeName(name))
				require.NoError(t, err)
				require.Equal(t, tftypes.NewValue(tftypes.String, expected), value)
			}
			// Destroy without refreshing first remains available for a cluster rejected during creation.
			destroyed, err := server.ApplyResourceChange(ctx, &tfprotov6.ApplyResourceChangeRequest{
				TypeName: "mongodbatlas_advanced_cluster", PriorState: created.NewState, PlannedState: empty, Config: empty,
			})
			require.NoError(t, err)
			require.Empty(t, destroyed.Diagnostics)
			destroyedState, err := destroyed.NewState.Unmarshal(typ)
			require.NoError(t, err)
			require.True(t, destroyedState.IsNull())
			require.Equal(t, []string{"POST", "GET", "DELETE", "GET"}, httpMock.methods)
		})
	}
}

// TestInfiniteClusterImportRejectsUnsupportedTopology covers the read path: an imported cluster reveals
// its topology only through the Atlas response, where effectiveDatabaseEdition is authoritative.
func TestInfiniteClusterImportRejectsUnsupportedTopology(t *testing.T) {
	for _, clusterType := range []string{"SHARDED", "GEOSHARDED", "REPLICASET"} {
		for _, effectiveEdition := range []string{"INFINITE", "CORE"} {
			t.Run(clusterType+"/"+effectiveEdition, func(t *testing.T) {
				for _, name := range []string{"ASSUME_ROLE_ARN", "TF_VAR_ASSUME_ROLE_ARN", "MONGODB_ATLAS_CLIENT_ID", "MONGODB_ATLAS_CLIENT_SECRET"} {
					t.Setenv(name, "")
				}
				var expected *regexp.Regexp
				if effectiveEdition == "INFINITE" && clusterType != "REPLICASET" {
					expected = regexp.MustCompile("Unsupported INFINITE cluster type")
				}
				httpMock := &infiniteImportHTTPMock{clusterType: clusterType, effectiveEdition: effectiveEdition}
				resource.UnitTest(t, resource.TestCase{
					PreCheck:                 func() { require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach()) },
					ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, httpMock),
					Steps: []resource.TestStep{{
						Config: topologyImportConfig(clusterType), ResourceName: "mongodbatlas_advanced_cluster.test",
						ImportState: true, ImportStateId: "111111111111111111111111-example", ExpectError: expected,
					}},
				})
			})
		}
	}
}

func topologyImportConfig(clusterType string) string {
	return fmt.Sprintf(`
provider "mongodbatlas" {
  public_key = "test-public-key"
  private_key = "test-private-key"
  base_url = "https://atlas.invalid/"
}
resource "mongodbatlas_advanced_cluster" "test" {
  name = "example"
  project_id = "111111111111111111111111"
  cluster_type = %q
  replication_specs = [{
    region_configs = [{
      provider_name = "AWS"
      region_name = "US_EAST_1"
      priority = 7
      electable_specs = { instance_size = "M10", node_count = 3 }
    }]
  }]
}
`, clusterType)
}

type infiniteImportHTTPMock struct {
	clusterType, effectiveEdition string
}

func (m *infiniteImportHTTPMock) ModifyHTTPClient(client *http.Client) error {
	client.Transport = m
	return nil
}

func (m *infiniteImportHTTPMock) ResetHTTPClient(*http.Client) {}

func (m *infiniteImportHTTPMock) RoundTrip(req *http.Request) (*http.Response, error) {
	var cluster map[string]any
	if err := json.Unmarshal([]byte(infiniteTopologyClusterResponse), &cluster); err != nil {
		return nil, err
	}
	cluster["clusterType"], cluster["effectiveDatabaseEdition"] = m.clusterType, m.effectiveEdition
	body, err := json.Marshal(cluster)
	if err != nil {
		return nil, err
	}
	if req.Method != http.MethodGet {
		return nil, fmt.Errorf("import must not write to Atlas: %s %s", req.Method, req.URL.Path)
	}
	switch {
	case strings.HasSuffix(req.URL.Path, "/clusters/example"):
	case strings.Contains(req.URL.Path, "/containers"):
		body = []byte(`{"results": [{"id": "666666666666666666666666", "providerName": "AWS", "regionName": "US_EAST_1"}], "totalCount": 1}`)
	default:
		body = []byte(`{}`)
	}
	return &http.Response{
		StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}},
		Body: io.NopCloser(strings.NewReader(string(body))), Request: req,
	}, nil
}

func topologyValidationAttributes(clusterType, edition any) map[string]any {
	return map[string]any{
		"name": "example", "project_id": dummyProjectID, "cluster_type": clusterType, "database_edition": edition,
		"replication_specs": []any{map[string]any{"region_configs": []any{map[string]any{
			"provider_name": "AWS", "region_name": "US_EAST_1", "priority": int64(7),
			"electable_specs": map[string]any{"instance_size": "M10", "node_count": int64(3)},
		}}}},
	}
}

func topologyValidationConfig(clusterType, edition string) string {
	editionAttribute := ""
	if edition != "" {
		editionAttribute = fmt.Sprintf("database_edition = %q", edition)
	}
	return fmt.Sprintf(`
provider "mongodbatlas" {
  public_key = "test-public-key"
  private_key = "test-private-key"
  base_url = "https://atlas.invalid/"
}
resource "mongodbatlas_advanced_cluster" "test" {
  name = "example"
  project_id = %q
  cluster_type = %q
  %s
  replication_specs = [{
    region_configs = [{
      provider_name = "AWS"
      region_name = "US_EAST_1"
      priority = 7
      electable_specs = { instance_size = "M10", node_count = 3 }
    }]
  }]
}
`, dummyProjectID, clusterType, editionAttribute)
}

// topologyTestValue fills omitted attributes with typed nulls so the fixture follows the current resource schema.
func topologyTestValue(typ tftypes.Type, value any) tftypes.Value {
	if value == nil || value == tftypes.UnknownValue {
		return tftypes.NewValue(typ, value)
	}
	switch shape := typ.(type) {
	case tftypes.Object:
		input := value.(map[string]any)
		out := map[string]tftypes.Value{}
		for name, childType := range shape.AttributeTypes {
			out[name] = topologyTestValue(childType, input[name])
		}
		return tftypes.NewValue(typ, out)
	case tftypes.List:
		out := []tftypes.Value{}
		for _, item := range value.([]any) {
			out = append(out, topologyTestValue(shape.ElementType, item))
		}
		return tftypes.NewValue(typ, out)
	default:
		return tftypes.NewValue(typ, value)
	}
}

type noTopologyRequests struct {
	calls atomic.Int64
}

func (m *noTopologyRequests) ModifyHTTPClient(client *http.Client) error {
	client.Transport = m
	return nil
}

func (m *noTopologyRequests) ResetHTTPClient(*http.Client) {}

func (m *noTopologyRequests) RoundTrip(req *http.Request) (*http.Response, error) {
	m.calls.Add(1)
	return nil, fmt.Errorf("unexpected Atlas request: %s %s", req.Method, req.URL.Path)
}

type unexpectedInfiniteCreateHTTPMock struct {
	clusterType string
	methods     []string
	deleted     bool
}

func (m *unexpectedInfiniteCreateHTTPMock) ModifyHTTPClient(client *http.Client) error {
	client.Transport = m
	return nil
}

func (m *unexpectedInfiniteCreateHTTPMock) ResetHTTPClient(*http.Client) {}

func (m *unexpectedInfiniteCreateHTTPMock) RoundTrip(req *http.Request) (*http.Response, error) {
	const clustersPath = "/api/atlas/v2/groups/111111111111111111111111/clusters"
	const clusterPath = clustersPath + "/example"
	m.methods = append(m.methods, req.Method)
	status := http.StatusOK
	var cluster map[string]any
	if err := json.Unmarshal([]byte(infiniteTopologyClusterResponse), &cluster); err != nil {
		return nil, err
	}
	cluster["clusterType"], cluster["effectiveDatabaseEdition"] = m.clusterType, "INFINITE"
	body, err := json.Marshal(cluster)
	if err != nil {
		return nil, err
	}
	switch {
	case req.Method == http.MethodPost && req.URL.Path == clustersPath:
		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			return nil, err
		}
		if payload["databaseEdition"] != nil || payload["paused"] != nil {
			return nil, fmt.Errorf("create must omit databaseEdition and paused")
		}
	case req.Method == http.MethodGet && req.URL.Path == clusterPath:
		if m.deleted {
			status, body = http.StatusNotFound, []byte(`{"errorCode":"CLUSTER_NOT_FOUND","error":404}`)
		}
	case req.Method == http.MethodDelete && req.URL.Path == clusterPath:
		m.deleted = true
		body = []byte(`{}`)
	default:
		return nil, fmt.Errorf("unexpected mocked Atlas request: %s %s", req.Method, req.URL.Path)
	}
	return &http.Response{
		StatusCode: status, Header: http.Header{"Content-Type": {"application/json"}},
		Body: io.NopCloser(strings.NewReader(string(body))), Request: req,
	}, nil
}

// infiniteTopologyClusterResponse omits databaseEdition: Atlas reports the edition it selected in effectiveDatabaseEdition.
const infiniteTopologyClusterResponse = `{
  "id": "333333333333333333333333", "groupId": "111111111111111111111111", "name": "example",
  "clusterType": "REPLICASET", "stateName": "IDLE",
  "mongoDBMajorVersion": "8.0", "mongoDBVersion": "8.0.5", "versionReleaseSystem": "LTS",
  "rootCertType": "ISRGROOTX1", "encryptionAtRestProvider": "NONE", "replicaSetScalingStrategy": "SEQUENTIAL",
  "replicationSpecs": [{
    "id": "444444444444444444444444", "zoneId": "555555555555555555555555", "zoneName": "Zone 1",
    "regionConfigs": [{
      "providerName": "AWS", "regionName": "US_EAST_1", "priority": 7,
      "electableSpecs": {"instanceSize": "M10", "nodeCount": 3, "diskSizeGB": 10, "diskIOPS": 3000, "ebsVolumeType": "STANDARD"}
    }]
  }]
}`

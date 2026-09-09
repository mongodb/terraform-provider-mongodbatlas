package advancedcluster_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
)

func TestPlanAutoScalingStorageConfigUnknown(t *testing.T) {
	ctx := t.Context()
	var schemaResponse frameworkresource.SchemaResponse
	advancedcluster.Resource().Schema(ctx, frameworkresource.SchemaRequest{}, &schemaResponse)
	typ := schemaResponse.Schema.Type().TerraformType(ctx)
	model := func(autoScaling any) *tfprotov6.DynamicValue {
		value := planTestValue(typ, map[string]any{
			"name": "example", "project_id": "111111111111111111111111", "cluster_type": "REPLICASET", "database_edition": "INFINITE",
			"replication_specs": []any{map[string]any{"region_configs": []any{map[string]any{
				"provider_name": "AWS", "region_name": "US_EAST_1", "priority": int64(7),
				"electable_specs": map[string]any{"instance_size": "M10", "node_count": int64(2)},
				"auto_scaling":    autoScaling,
			}}}},
		})
		result, err := tfprotov6.NewDynamicValue(typ, value)
		require.NoError(t, err)
		return &result
	}
	settings := func(storage any) map[string]any {
		return map[string]any{
			"compute_enabled": true, "compute_max_instance_size": "M20", "compute_scale_down_enabled": false,
			"storage_config": storage,
		}
	}
	autoScalingPath := tftypes.NewAttributePath().WithAttributeName("replication_specs").WithElementKeyInt(0).
		WithAttributeName("region_configs").WithElementKeyInt(0).WithAttributeName("auto_scaling")
	storagePath := autoScalingPath.WithAttributeName("storage_config")
	limitPath := storagePath.WithAttributeName("shard_size_limit_gb")
	for _, priorStorage := range []bool{false, true} {
		var storage any
		if priorStorage {
			storage = map[string]any{"shard_size_limit_gb": int64(1024)}
		}
		prior := model(settings(storage))
		for name, tc := range map[string]struct {
			config      any
			unknownPath *tftypes.AttributePath
		}{
			"shard limit": {
				config: settings(map[string]any{"shard_size_limit_gb": tftypes.UnknownValue}), unknownPath: limitPath,
			},
			"storage config": {
				config: settings(tftypes.UnknownValue), unknownPath: storagePath,
			},
			"auto scaling": {
				config: tftypes.UnknownValue, unknownPath: autoScalingPath,
			},
		} {
			t.Run(fmt.Sprintf("prior_storage=%t/%s", priorStorage, name), func(t *testing.T) {
				server, err := acc.TestAccProviderV6Factories["mongodbatlas"]()
				require.NoError(t, err)
				planValue := func(config *tfprotov6.DynamicValue, path *tftypes.AttributePath) tftypes.Value {
					result, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
						TypeName: "mongodbatlas_advanced_cluster", PriorState: prior, ProposedNewState: config, Config: config,
					})
					require.NoError(t, err)
					require.Empty(t, result.Diagnostics)
					plan, err := result.PlannedState.Unmarshal(typ)
					require.NoError(t, err)
					value, _, err := tftypes.WalkAttributePath(plan, path)
					require.NoError(t, err)
					return value.(tftypes.Value)
				}
				initial := planValue(model(tc.config), tc.unknownPath)
				require.False(t, initial.IsKnown(), "configured unknown storage values must not be replaced by prior state")
				resolved := planValue(model(settings(map[string]any{"shard_size_limit_gb": int64(2048)})), limitPath)
				require.Equal(t, tftypes.NewValue(tftypes.Number, int64(2048)), resolved, "the plan must accept the value resolved during apply")
			})
		}
	}
}

func TestAutoScalingStorageConfigUnknownLifecycle(t *testing.T) {
	t.Setenv("ASSUME_ROLE_ARN", "")
	t.Setenv("TF_VAR_ASSUME_ROLE_ARN", "")
	t.Setenv("MONGODB_ATLAS_CLIENT_ID", "")
	t.Setenv("MONGODB_ATLAS_CLIENT_SECRET", "")
	mock := &storageConfigHTTPMock{}
	config := func(limit int) string {
		return fmt.Sprintf(`
provider "mongodbatlas" {
  public_key = "test-public-key"
  private_key = "test-private-key"
  base_url = "https://atlas.invalid/"
}
resource "terraform_data" "limit" { input = %[1]d }
resource "mongodbatlas_advanced_cluster" "test" {
  name = "example"
  project_id = "111111111111111111111111"
  cluster_type = "REPLICASET"
  database_edition = "INFINITE"
  replication_specs = [{
    region_configs = [{
      provider_name = "AWS"
      region_name = "US_EAST_1"
      priority = 7
      electable_specs = { instance_size = "M10", node_count = 2 }
      auto_scaling = {
        compute_enabled = true
        compute_max_instance_size = "M20"
        compute_scale_down_enabled = false
        storage_config = { shard_size_limit_gb = terraform_data.limit.output }
      }
    }]
  }]
}
`, limit)
	}
	const limitAttribute = "replication_specs.0.region_configs.0.auto_scaling.storage_config.shard_size_limit_gb"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach()) },
		ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, mock),
		Steps: []resource.TestStep{
			{Config: config(1024), Check: resource.TestCheckResourceAttr("mongodbatlas_advanced_cluster.test", limitAttribute, "1024")},
			{Config: config(2048), Check: resource.TestCheckResourceAttr("mongodbatlas_advanced_cluster.test", limitAttribute, "2048")},
		},
	})
}

type storageConfigHTTPMock struct {
	mu      sync.Mutex
	limit   int
	deleted bool
}

func (m *storageConfigHTTPMock) ModifyHTTPClient(client *http.Client) error {
	client.Transport = m
	return nil
}

func (m *storageConfigHTTPMock) ResetHTTPClient(*http.Client) {}

func (m *storageConfigHTTPMock) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	const groupPath = "/api/atlas/v2/groups/111111111111111111111111"
	const clusterPath = groupPath + "/clusters/example"
	status := http.StatusOK
	body := "{}"
	switch {
	case req.Method == http.MethodPost && req.URL.Path == groupPath+"/clusters",
		req.Method == http.MethodPatch && req.URL.Path == clusterPath:
		var cluster admin.ClusterDescription20240805
		if err := json.NewDecoder(req.Body).Decode(&cluster); err != nil {
			return nil, err
		}
		specs := cluster.GetReplicationSpecs()
		if len(specs) != 1 || len(specs[0].GetRegionConfigs()) != 1 {
			return nil, fmt.Errorf("expected a single replication spec and region in %s request", req.Method)
		}
		autoScaling := specs[0].GetRegionConfigs()[0].GetAutoScaling()
		storage := autoScaling.GetStorageConfig()
		if !storage.HasShardSizeLimitGB() {
			return nil, fmt.Errorf("expected shardSizeLimitGB in %s request", req.Method)
		}
		m.limit = storage.GetShardSizeLimitGB()
		body = fmt.Sprintf(storageConfigClusterResponse, m.limit)
	case req.Method == http.MethodGet && req.URL.Path == clusterPath:
		body = fmt.Sprintf(storageConfigClusterResponse, m.limit)
		if m.deleted {
			status, body = http.StatusNotFound, `{"errorCode":"CLUSTER_NOT_FOUND","error":404}`
		}
	case req.Method == http.MethodGet && req.URL.Path == groupPath+"/containers":
		body = `{"results":[{"id":"222222222222222222222222","providerName":"AWS","regionName":"US_EAST_1"}],"totalCount":1}`
	case req.Method == http.MethodGet && req.URL.Path == clusterPath+"/processArgs":
	case req.Method == http.MethodDelete && req.URL.Path == clusterPath:
		m.deleted = true
	default:
		return nil, fmt.Errorf("unexpected mocked Atlas request: %s %s", req.Method, req.URL.Path)
	}
	return &http.Response{
		StatusCode: status, Header: http.Header{"Content-Type": {"application/json"}},
		Body: io.NopCloser(strings.NewReader(body)), Request: req,
	}, nil
}

const storageConfigClusterResponse = `{
  "id": "333333333333333333333333", "groupId": "111111111111111111111111", "name": "example",
  "clusterType": "REPLICASET", "databaseEdition": "INFINITE", "stateName": "IDLE",
  "mongoDBMajorVersion": "8.0", "mongoDBVersion": "8.0.5", "versionReleaseSystem": "LTS",
  "rootCertType": "ISRGROOTX1", "encryptionAtRestProvider": "NONE", "replicaSetScalingStrategy": "SEQUENTIAL",
  "replicationSpecs": [{
    "id": "444444444444444444444444", "zoneId": "555555555555555555555555", "zoneName": "Zone 1",
    "regionConfigs": [{
      "providerName": "AWS", "regionName": "US_EAST_1", "priority": 7,
      "electableSpecs": {"instanceSize": "M10", "nodeCount": 2, "diskSizeGB": 10, "diskIOPS": 3000, "ebsVolumeType": "STANDARD"},
      "autoScaling": {
        "compute": {"enabled": true, "maxInstanceSize": "M20", "minInstanceSize": "M10", "scaleDownEnabled": false},
        "diskGB": {"enabled": false}, "storageConfig": {"shardSizeLimitGB": %d}
      }
    }]
  }]
}`

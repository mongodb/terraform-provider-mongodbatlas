package advancedcluster_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
)

func TestAutoScalingStorageConfigImportLifecycle(t *testing.T) {
	const configuredLimit = 1024
	// The import kind and the requested edition are independent, so pair them instead of running the full matrix.
	for _, tc := range []struct{ importBlock, requestedEdition bool }{{importBlock: false, requestedEdition: true}, {importBlock: true, requestedEdition: false}} {
		importBlock, requestedEdition := tc.importBlock, tc.requestedEdition
		t.Run(fmt.Sprintf("block=%t/requested_edition=%t", importBlock, requestedEdition), func(t *testing.T) {
			storageImportCredentials(t)
			mock := &storageImportHTTPMock{clusterType: "REPLICASET", requestedEdition: requestedEdition, effectiveEdition: "INFINITE", limit: configuredLimit}
			configForLimit := func(limit int, omitAutoScaling bool) string {
				config := storageImportConfig(requestedEdition, limit, omitAutoScaling)
				if !requestedEdition {
					config = strings.ReplaceAll(config, `cluster_type = "REPLICASET"`, "cluster_type = \"REPLICASET\"\n  use_effective_fields = true")
				}
				return config
			}
			config := configForLimit(configuredLimit, false)
			check := storageImportStateCheck(configuredLimit)
			first := resource.TestStep{
				Config: config, ResourceName: "mongodbatlas_advanced_cluster.test", ImportState: true,
				ImportStateId: storageImportID, ImportStatePersist: true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					require.Len(t, states, 1)
					require.Equal(t, "1024", states[0].Attributes[storageImportLimitAttribute])
					return nil
				},
			}
			if importBlock {
				first = resource.TestStep{Config: config + storageImportBlock, Check: check}
			}
			resource.UnitTest(t, resource.TestCase{
				PreCheck:                 func() { require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach()) },
				ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, mock),
				Steps: []resource.TestStep{
					first,
					{Config: config, Check: check},
					{Config: configForLimit(2048, false), Check: storageImportStateCheck(2048)},
					{Config: configForLimit(0, importBlock), Check: storageImportStateCheck(0)},
				},
			})
			expectedLimits := []int{2048, 0}
			if !requestedEdition {
				// Enabling effective fields after import resends the requested replication specs once.
				expectedLimits = append([]int{configuredLimit}, expectedLimits...)
				require.True(t, mock.usedEffectiveFields)
			}
			require.Equal(t, expectedLimits, mock.updatedLimits)
		})
	}
}

// TestAutoScalingStorageConfigImportReadsMissingLimit covers the import-time read of an
// unconfigured shard limit; the full update lifecycle above always starts from a configured one.
func TestAutoScalingStorageConfigImportReadsMissingLimit(t *testing.T) {
	storageImportCredentials(t)
	mock := &storageImportHTTPMock{clusterType: "REPLICASET", requestedEdition: true, effectiveEdition: "INFINITE"}
	config := storageImportConfig(true, 0, false)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach()) },
		ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, mock),
		Steps: []resource.TestStep{{
			Config: config, ResourceName: "mongodbatlas_advanced_cluster.test", ImportState: true,
			ImportStateId: storageImportID,
			ImportStateCheck: func(states []*terraform.InstanceState) error {
				require.Len(t, states, 1)
				require.Empty(t, states[0].Attributes[storageImportLimitAttribute])
				return nil
			},
		}},
	})
}

func TestInfiniteClusterImportWithoutRequestedEditionOrShardLimit(t *testing.T) {
	for _, explicitFalse := range []bool{false, true} {
		t.Run(fmt.Sprintf("explicit_false=%t", explicitFalse), func(t *testing.T) {
			storageImportCredentials(t)
			mock := &storageImportHTTPMock{clusterType: "REPLICASET", effectiveEdition: "INFINITE", omitDisk: !explicitFalse, explicitFalse: explicitFalse}
			config := storageImportConfig(false, 0, false)
			if explicitFalse {
				config = strings.ReplaceAll(config, "compute_enabled = true", "disk_gb_enabled = false\n        compute_enabled = true")
			}
			resource.UnitTest(t, resource.TestCase{
				PreCheck:                 func() { require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach()) },
				ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, mock),
				Steps: []resource.TestStep{
					{
						Config: config, ResourceName: "mongodbatlas_advanced_cluster.test", ImportState: true,
						ImportStateId: storageImportID, ImportStatePersist: true,
					},
					{
						Config: strings.ReplaceAll(config, "priority = 7", "priority = 6"),
						Check:  resource.TestCheckResourceAttr("mongodbatlas_advanced_cluster.test", "replication_specs.0.region_configs.0.priority", "6"),
					},
				},
			})
		})
	}
}

func TestAutoScalingStorageConfigImportClearsUnconfiguredLimit(t *testing.T) {
	// importBlock vs ResourceName-based import is already cross-tested by the lifecycle
	// test above; clearing an unconfigured limit is orthogonal to that choice.
	storageImportCredentials(t)
	mock := &storageImportHTTPMock{clusterType: "REPLICASET", effectiveEdition: "INFINITE", limit: 1024}
	config := storageImportConfig(false, 0, true)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach()) },
		ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, mock),
		Steps: []resource.TestStep{
			{
				Config: config, ResourceName: "mongodbatlas_advanced_cluster.test", ImportState: true,
				ImportStateId: storageImportID, ImportStatePersist: true,
			},
			{Config: config, Check: storageImportStateCheck(0)},
		},
	})
	require.Equal(t, []int{0}, mock.updatedLimits)
}

func TestCoreClusterImportSupportsAllTopologies(t *testing.T) {
	// GEOSHARDED is omitted: it shares the SHARDED import path.
	for _, clusterType := range []string{"REPLICASET", "SHARDED"} {
		t.Run(clusterType, func(t *testing.T) {
			storageImportCredentials(t)
			mock := &storageImportHTTPMock{clusterType: clusterType, effectiveEdition: "CORE"}
			config := strings.ReplaceAll(storageImportConfig(false, 0, false), `"REPLICASET"`, fmt.Sprintf("%q", clusterType))
			resource.UnitTest(t, resource.TestCase{
				PreCheck:                 func() { require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach()) },
				ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, mock),
				Steps: []resource.TestStep{
					{
						Config: config, ResourceName: "mongodbatlas_advanced_cluster.test", ImportState: true,
						ImportStateId: storageImportID, ImportStatePersist: true,
					},
					{
						Config: config,
						Check:  resource.TestCheckResourceAttr("mongodbatlas_advanced_cluster.test", "cluster_type", clusterType),
					},
				},
			})
			require.Empty(t, mock.updatedLimits)
		})
	}
}

func storageImportCredentials(t *testing.T) {
	t.Helper()
	t.Setenv("ASSUME_ROLE_ARN", "")
	t.Setenv("TF_VAR_ASSUME_ROLE_ARN", "")
	t.Setenv("MONGODB_ATLAS_CLIENT_ID", "")
	t.Setenv("MONGODB_ATLAS_CLIENT_SECRET", "")
}

func storageImportStateCheck(limit int) resource.TestCheckFunc {
	if limit == 0 {
		return resource.TestCheckNoResourceAttr("mongodbatlas_advanced_cluster.test", storageImportLimitAttribute)
	}
	return resource.TestCheckResourceAttr("mongodbatlas_advanced_cluster.test", storageImportLimitAttribute, fmt.Sprint(limit))
}

func storageImportConfig(requestedEdition bool, limit int, omitAutoScaling bool) string {
	edition := ""
	if requestedEdition {
		edition = `database_edition = "INFINITE"`
	}
	storage := ""
	if limit != 0 {
		storage = fmt.Sprintf("storage_config = { shard_size_limit_gb = %d }", limit)
	}
	autoScaling := ""
	if !omitAutoScaling {
		autoScaling = fmt.Sprintf(`auto_scaling = {
        compute_enabled = true
        compute_max_instance_size = "M20"
        compute_scale_down_enabled = false
        %s
      }`, storage)
	}
	return fmt.Sprintf(`
provider "mongodbatlas" {
  public_key = "test-public-key"
  private_key = "test-private-key"
  base_url = "https://atlas.invalid/"
}
resource "mongodbatlas_advanced_cluster" "test" {
  name = "example-with-hyphens"
  project_id = "111111111111111111111111"
  cluster_type = "REPLICASET"
  %[1]s
  replication_specs = [{
    region_configs = [{
      provider_name = "AWS"
      region_name = "US_EAST_1"
      priority = 7
      electable_specs = { instance_size = "M10", node_count = 2 }
      %[2]s
    }]
  }]
}
`, edition, autoScaling)
}

type storageImportHTTPMock struct {
	clusterType         string
	effectiveEdition    string
	updatedLimits       []int
	mu                  sync.Mutex
	limit               int
	priority            int
	requestedEdition    bool
	omitDisk            bool
	explicitFalse       bool
	usedEffectiveFields bool
	deleted             bool
}

func (m *storageImportHTTPMock) ModifyHTTPClient(client *http.Client) error {
	client.Transport = m
	return nil
}

func (m *storageImportHTTPMock) ResetHTTPClient(*http.Client) {}

func (m *storageImportHTTPMock) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	const groupPath = "/api/atlas/v2/groups/111111111111111111111111"
	const clusterPath = groupPath + "/clusters/example-with-hyphens"
	status := http.StatusOK
	body := "{}"
	switch {
	case req.Method == http.MethodGet && req.URL.Path == clusterPath,
		req.Method == http.MethodPatch && req.URL.Path == clusterPath:
		if req.Header.Get("Use-Effective-Instance-Fields") == "true" {
			m.usedEffectiveFields = true
		}
		if req.Method == http.MethodPatch {
			var patch admin.ClusterDescription20240805
			if err := json.NewDecoder(req.Body).Decode(&patch); err != nil {
				return nil, err
			}
			specs := patch.GetReplicationSpecs()
			if len(specs) != 1 || len(specs[0].GetRegionConfigs()) != 1 {
				return nil, fmt.Errorf("import lifecycle update must include one replication spec and region")
			}
			region := specs[0].GetRegionConfigs()[0]
			autoScaling := region.GetAutoScaling()
			if autoScaling.DiskGB != nil && !autoScaling.DiskGB.HasEnabled() {
				return nil, fmt.Errorf("INFINITE requests must not include converter-generated diskGB: {}")
			}
			if m.explicitFalse && (autoScaling.DiskGB == nil || !autoScaling.DiskGB.HasEnabled() || autoScaling.DiskGB.GetEnabled()) {
				return nil, fmt.Errorf("configured diskGB.enabled=false must remain in the update request")
			}
			if m.explicitFalse && (autoScaling.Compute == nil || !autoScaling.Compute.HasScaleDownEnabled() || autoScaling.Compute.GetScaleDownEnabled()) {
				return nil, fmt.Errorf("configured compute.scaleDownEnabled=false must remain in the update request")
			}
			storage := autoScaling.GetStorageConfig()
			m.limit = storage.GetShardSizeLimitGB()
			m.priority = region.GetPriority()
			m.updatedLimits = append(m.updatedLimits, m.limit)
		}
		var cluster map[string]any
		if err := json.Unmarshal(fmt.Appendf(nil, storageConfigClusterResponse, m.limit), &cluster); err != nil {
			return nil, err
		}
		cluster["clusterType"] = m.clusterType
		cluster["name"] = "example-with-hyphens"
		cluster["effectiveDatabaseEdition"] = m.effectiveEdition
		if !m.requestedEdition {
			delete(cluster, "databaseEdition")
		}
		region := cluster["replicationSpecs"].([]any)[0].(map[string]any)["regionConfigs"].([]any)[0].(map[string]any)
		if m.priority != 0 {
			region["priority"] = m.priority
		}
		autoScaling := region["autoScaling"].(map[string]any)
		if m.omitDisk {
			delete(autoScaling, "diskGB")
		}
		if m.limit == 0 {
			if m.requestedEdition {
				delete(autoScaling, "storageConfig")
			} else {
				autoScaling["storageConfig"] = map[string]any{}
			}
		}
		encoded, err := json.Marshal(cluster)
		if err != nil {
			return nil, err
		}
		body = string(encoded)
		if m.deleted {
			status, body = http.StatusNotFound, `{"errorCode":"CLUSTER_NOT_FOUND","error":404}`
		}
	case req.Method == http.MethodGet && req.URL.Path == groupPath+"/containers":
		body = `{"results":[{"id":"222222222222222222222222","providerName":"AWS","regionName":"US_EAST_1"}],"totalCount":1}`
	case req.Method == http.MethodGet && req.URL.Path == clusterPath+"/processArgs":
	case req.Method == http.MethodDelete && req.URL.Path == clusterPath:
		m.deleted = true
	default:
		return nil, fmt.Errorf("unexpected mocked import request: %s %s", req.Method, req.URL.Path)
	}
	return &http.Response{
		StatusCode: status, Header: http.Header{"Content-Type": {"application/json"}},
		Body: io.NopCloser(strings.NewReader(body)), Request: req,
	}, nil
}

const storageImportLimitAttribute = "replication_specs.0.region_configs.0.auto_scaling.storage_config.shard_size_limit_gb"

const storageImportID = "111111111111111111111111-example-with-hyphens"

const storageImportBlock = `
import {
  to = mongodbatlas_advanced_cluster.test
  id = "111111111111111111111111-example-with-hyphens"
}
`

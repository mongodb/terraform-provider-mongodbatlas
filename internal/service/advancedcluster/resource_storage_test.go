package advancedcluster_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312026/admin"
)

// storageConfigCredentials clears ambient credentials so the mocked provider never reaches real Atlas.
func storageConfigCredentials(t *testing.T) {
	t.Helper()
	t.Setenv("ASSUME_ROLE_ARN", "")
	t.Setenv("TF_VAR_ASSUME_ROLE_ARN", "")
	t.Setenv("MONGODB_ATLAS_CLIENT_ID", "")
	t.Setenv("MONGODB_ATLAS_CLIENT_SECRET", "")
}

const storageConfigLimitAttribute = "replication_specs.0.region_configs.0.auto_scaling.storage_config.shard_size_limit_gb"

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

func TestPlanAutoScalingStorageConfigUnknown(t *testing.T) {
	ctx := t.Context()
	_, typ := clusterSchema(ctx, t)
	model := func(autoScaling any) *tfprotov6.DynamicValue {
		return clusterDynamic(t, typ, map[string]any{
			"name": "example", "project_id": "111111111111111111111111", "cluster_type": "REPLICASET", "database_edition": "INFINITE",
			"replication_specs": []any{map[string]any{"region_configs": []any{map[string]any{
				"provider_name": "AWS", "region_name": "US_EAST_1", "priority": int64(7),
				"electable_specs": map[string]any{"instance_size": "M10", "node_count": int64(2)},
				"auto_scaling":    autoScaling,
			}}}},
		})
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
	storageConfigCredentials(t)
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
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach()) },
		ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, mock),
		Steps: []resource.TestStep{
			{Config: config(1024), Check: resource.TestCheckResourceAttr("mongodbatlas_advanced_cluster.test", storageConfigLimitAttribute, "1024")},
			{Config: config(2048), Check: resource.TestCheckResourceAttr("mongodbatlas_advanced_cluster.test", storageConfigLimitAttribute, "2048")},
		},
	})
}

// storageConfigHTTPMock serves a cluster whose configured shard limit tracks the last value written.
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

func TestAutoScalingStorageConfigImportLifecycle(t *testing.T) {
	const configuredLimit = 1024
	// The effective variant runs the import lifecycle with use_effective_fields on an INFINITE cluster
	// without the requested edition, expecting the Use-Effective-Instance-Fields header on updates.
	for _, effective := range []bool{false, true} {
		for _, importBlock := range []bool{false, true} {
			t.Run(fmt.Sprintf("effective=%t/block=%t", effective, importBlock), func(t *testing.T) {
				storageConfigCredentials(t)
				mock := &storageImportHTTPMock{clusterType: "REPLICASET", requestedEdition: !effective, effectiveEdition: "INFINITE", limit: configuredLimit, expectEffectiveFields: effective}
				configForLimit := func(limit int, omitAutoScaling bool) string {
					return storageImportConfig(!effective, limit, omitAutoScaling, effective)
				}
				config := configForLimit(configuredLimit, false)
				check := storageImportStateCheck(configuredLimit)
				first := resource.TestStep{
					Config: config, ResourceName: "mongodbatlas_advanced_cluster.test", ImportState: true,
					ImportStateId: storageImportID, ImportStatePersist: true,
					ImportStateCheck: func(states []*terraform.InstanceState) error {
						require.Len(t, states, 1)
						require.Equal(t, "1024", states[0].Attributes[storageConfigLimitAttribute])
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
				// The flag-toggle apply skips the copy-from-state optimization, re-sending the configured
				// replication specs (with the current limit) once before the actual updates.
				expectedLimits := []int{2048, 0}
				if effective {
					expectedLimits = []int{1024, 2048, 0}
				}
				require.Equal(t, expectedLimits, mock.updatedLimits)
			})
		}
	}
}

// TestAutoScalingStorageConfigImportReadsMissingLimit covers the import-time read of an
// unconfigured shard limit; the full update lifecycle above always starts from a configured one.
func TestAutoScalingStorageConfigImportReadsMissingLimit(t *testing.T) {
	storageConfigCredentials(t)
	mock := &storageImportHTTPMock{clusterType: "REPLICASET", requestedEdition: true, effectiveEdition: "INFINITE"}
	config := storageImportConfig(true, 0, false, false)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach()) },
		ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, mock),
		Steps: []resource.TestStep{{
			Config: config, ResourceName: "mongodbatlas_advanced_cluster.test", ImportState: true,
			ImportStateId: storageImportID,
			ImportStateCheck: func(states []*terraform.InstanceState) error {
				require.Len(t, states, 1)
				require.Empty(t, states[0].Attributes[storageConfigLimitAttribute])
				return nil
			},
		}},
	})
}

func TestInfiniteClusterImportWithoutRequestedEditionOrShardLimit(t *testing.T) {
	for _, explicitFalse := range []bool{false, true} {
		t.Run(fmt.Sprintf("explicit_false=%t", explicitFalse), func(t *testing.T) {
			storageConfigCredentials(t)
			mock := &storageImportHTTPMock{clusterType: "REPLICASET", effectiveEdition: "INFINITE", omitDisk: !explicitFalse, explicitFalse: explicitFalse}
			config := storageImportConfig(false, 0, false, false)
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
	storageConfigCredentials(t)
	mock := &storageImportHTTPMock{clusterType: "REPLICASET", effectiveEdition: "INFINITE", limit: 1024}
	config := storageImportConfig(false, 0, true, false)
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
			storageConfigCredentials(t)
			mock := &storageImportHTTPMock{clusterType: clusterType, effectiveEdition: "CORE"}
			config := strings.ReplaceAll(storageImportConfig(false, 0, false, false), `"REPLICASET"`, fmt.Sprintf("%q", clusterType))
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

func storageImportStateCheck(limit int) resource.TestCheckFunc {
	if limit == 0 {
		return resource.TestCheckNoResourceAttr("mongodbatlas_advanced_cluster.test", storageConfigLimitAttribute)
	}
	return resource.TestCheckResourceAttr("mongodbatlas_advanced_cluster.test", storageConfigLimitAttribute, fmt.Sprint(limit))
}

func storageImportConfig(requestedEdition bool, limit int, omitAutoScaling, useEffectiveFields bool) string {
	edition := ""
	if requestedEdition {
		edition = `database_edition = "INFINITE"`
	}
	effectiveFields := ""
	if useEffectiveFields {
		effectiveFields = "use_effective_fields = true"
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
  %[3]s
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
`, edition, autoScaling, effectiveFields)
}

// storageImportHTTPMock serves imported clusters across editions and topologies, tracking written limits.
type storageImportHTTPMock struct {
	clusterType           string
	effectiveEdition      string
	updatedLimits         []int
	mu                    sync.Mutex
	limit                 int
	priority              int
	requestedEdition      bool
	omitDisk              bool
	explicitFalse         bool
	expectEffectiveFields bool
	deleted               bool
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
		if req.Method == http.MethodPatch {
			if got := req.Header.Get("Use-Effective-Instance-Fields"); (got == "true") != m.expectEffectiveFields {
				return nil, fmt.Errorf("Use-Effective-Instance-Fields header = %q, expected effective fields %t", got, m.expectEffectiveFields)
			}
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

const storageImportID = "111111111111111111111111-example-with-hyphens"

const storageImportBlock = `
import {
  to = mongodbatlas_advanced_cluster.test
  id = "111111111111111111111111-example-with-hyphens"
}
`

// TestExplicitNullStorageConfig verifies that setting storageConfig to explicit null
// creates a detectable change in the PATCH payload.
func TestExplicitNullStorageConfig(t *testing.T) {
	// State: has storageConfig with shardSizeLimitGB
	stateReq := &admin.ClusterDescription20240805{
		ReplicationSpecs: &[]admin.ReplicationSpec20240805{
			{
				RegionConfigs: &[]admin.CloudRegionConfig20240805{
					{
						ProviderName: new("AWS"),
						RegionName:   new("US_EAST_1"),
						Priority:     new(7),
						AutoScaling: &admin.AdvancedAutoScalingSettings{
							Compute: &admin.AdvancedComputeAutoScaling{
								Enabled:          new(true),
								MaxInstanceSize:  new("M20"),
								MinInstanceSize:  new("M10"),
								ScaleDownEnabled: new(false),
							},
							StorageConfig: &admin.StorageConfig{
								ShardSizeLimitGB: new(1024),
							},
						},
					},
				},
			},
		},
	}

	// Plan: has storageConfig explicitly set to null
	planReq := &admin.ClusterDescription20240805{
		ReplicationSpecs: &[]admin.ReplicationSpec20240805{
			{
				RegionConfigs: &[]admin.CloudRegionConfig20240805{
					{
						ProviderName: new("AWS"),
						RegionName:   new("US_EAST_1"),
						Priority:     new(7),
						AutoScaling: &admin.AdvancedAutoScalingSettings{
							Compute: &admin.AdvancedComputeAutoScaling{
								Enabled:          new(true),
								MaxInstanceSize:  new("M20"),
								MinInstanceSize:  new("M10"),
								ScaleDownEnabled: new(false),
							},
						},
					},
				},
			},
		},
	}

	// Set storageConfig to explicit null in the plan
	planRegion := &(*planReq.ReplicationSpecs)[0].GetRegionConfigs()[0]
	planRegion.AutoScaling.SetStorageConfigNil()

	// Verify the plan serializes with explicit null
	planJSON, err := json.Marshal(planReq)
	require.NoError(t, err)
	require.Contains(t, string(planJSON), `"storageConfig":null`, "plan should serialize with explicit null")

	// Verify the diff detects a change (replacement, not removal)
	patchReq, err := update.PatchPayload(stateReq, planReq)
	require.NoError(t, err)
	require.NotNil(t, patchReq, "patch should not be nil when storageConfig is explicitly set to null")
	require.NotNil(t, patchReq.ReplicationSpecs, "patch should include replicationSpecs")

	// Verify the final serialized request
	patchJSON, err := json.Marshal(patchReq)
	require.NoError(t, err)

	// The PATCH should include replicationSpecs with autoScaling and compute settings preserved
	expected := `{"replicationSpecs":[{"regionConfigs":[{"priority":7,"providerName":"AWS","regionName":"US_EAST_1","autoScaling":{"compute":{"enabled":true,"maxInstanceSize":"M20","minInstanceSize":"M10","scaleDownEnabled":false}}}]}]}`
	require.JSONEq(t, expected, string(patchJSON), "PATCH should include replicationSpecs with compute settings preserved and storageConfig removed")
}

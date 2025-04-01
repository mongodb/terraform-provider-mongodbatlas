package advancedclustertpf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

var (
	repSpec0        = tfjsonpath.New("replication_specs").AtSliceIndex(0)
	repSpec1        = tfjsonpath.New("replication_specs").AtSliceIndex(1)
	regionConfig0_0 = repSpec0.AtMapKey("region_configs").AtSliceIndex(0)
	regionConfig0_1 = repSpec0.AtMapKey("region_configs").AtSliceIndex(1)
	regionConfig1_0 = repSpec1.AtMapKey("region_configs").AtSliceIndex(0)
	regionConfig1_1 = repSpec1.AtMapKey("region_configs").AtSliceIndex(1)
	advConfig       = tfjsonpath.New("advanced_configuration")
	mockConfig      = unit.MockConfigAdvancedClusterTPF
)

func autoScalingKnownValue(computeEnabled, diskEnabled, scaleDown bool, minInstanceSize, maxInstanceSize string) knownvalue.Check {
	return knownvalue.ObjectExact(map[string]knownvalue.Check{
		"compute_enabled":            knownvalue.Bool(computeEnabled),
		"disk_gb_enabled":            knownvalue.Bool(diskEnabled),
		"compute_scale_down_enabled": knownvalue.Bool(scaleDown),
		"compute_min_instance_size":  knownvalue.StringExact(minInstanceSize),
		"compute_max_instance_size":  knownvalue.StringExact(maxInstanceSize),
	})
}

func specInstanceSizeNodeCount(instanceSize string, nodeCount int) knownvalue.Check {
	return knownvalue.ObjectPartial(map[string]knownvalue.Check{
		"instance_size": knownvalue.StringExact(instanceSize),
		"node_count":    knownvalue.Int64Exact(int64(nodeCount)),
	})
}

func TestPlanChecksClusterTwoRepSpecsWithAutoScalingAndSpecs(t *testing.T) {
	var (
		baseConfig         = unit.NewMockPlanChecksConfig(t, &mockConfig, unit.ImportNameClusterTwoRepSpecsWithAutoScalingAndSpecs)
		resourceName       = baseConfig.ResourceName
		autoScalingEnabled = autoScalingKnownValue(true, true, true, "M10", "M30")
		testCases          = []unit.PlanCheckTest{
			{
				ConfigFilename: "main_removed_blocks_from_config_no_plan_changes.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
				},
			},
			{
				ConfigFilename: "main_removed_blocks_from_config_and_instance_change.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					// checks regionConfig0
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("read_only_specs"), specInstanceSizeNodeCount("M10", 2)),
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("electable_specs"), specInstanceSizeNodeCount("M10", 5)),
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("auto_scaling"), autoScalingEnabled),
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("analytics_auto_scaling"), autoScalingEnabled),
					plancheck.ExpectUnknownValue(resourceName, regionConfig0_0.AtMapKey("analytics_specs")), // analytics specs was defined in region_configs.0 but not in region_configs.1
					plancheck.ExpectUnknownValue(resourceName, repSpec0.AtMapKey("id")),

					// checks regionConfig1
					plancheck.ExpectKnownValue(resourceName, regionConfig1_0.AtMapKey("read_only_specs"), specInstanceSizeNodeCount("M20", 1)),
					plancheck.ExpectKnownValue(resourceName, regionConfig1_0.AtMapKey("electable_specs"), specInstanceSizeNodeCount("M20", 3)),
					plancheck.ExpectKnownValue(resourceName, regionConfig1_0.AtMapKey("auto_scaling"), autoScalingEnabled),
					plancheck.ExpectKnownValue(resourceName, regionConfig1_0.AtMapKey("analytics_auto_scaling"), autoScalingEnabled),
					plancheck.ExpectKnownValue(resourceName, regionConfig1_0.AtMapKey("analytics_specs"), knownvalue.NotNull()),
					plancheck.ExpectUnknownValue(resourceName, repSpec1.AtMapKey("id")),
				},
			},
		}
	)
	unit.RunPlanCheckTests(t, baseConfig, testCases)
}

func TestMockPlanChecks_ClusterReplicasetOneRegion(t *testing.T) {
	var (
		baseConfig   = unit.NewMockPlanChecksConfig(t, &mockConfig, unit.ImportNameClusterReplicasetOneRegion)
		resourceName = baseConfig.ResourceName
		testCases    = []unit.PlanCheckTest{
			{
				ConfigFilename: "main_mongo_db_major_version_changed.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("mongo_db_version")),
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("analytics_specs"), specInstanceSizeNodeCount("M10", 0)),
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("read_only_specs"), specInstanceSizeNodeCount("M10", 0)),
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("analytics_auto_scaling"), knownvalue.Null()),
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("auto_scaling"), autoScalingKnownValue(false, true, false, "", "")),
				},
			},
			{
				ConfigFilename: "main_backup_enabled.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					// should use state values inside replication_specs as no changes are made to replication_specs
					plancheck.ExpectKnownValue(resourceName, repSpec0.AtMapKey("id"), knownvalue.NotNull()),
					plancheck.ExpectKnownValue(resourceName, repSpec0.AtMapKey("zone_name"), knownvalue.NotNull()),
					plancheck.ExpectKnownValue(resourceName, repSpec0.AtMapKey("zone_id"), knownvalue.NotNull()),
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("electable_specs").AtMapKey("ebs_volume_type"), knownvalue.NotNull()),
				},
			},
			{
				ConfigFilename: "main_electable_disk_size_changed.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("disk_size_gb")),
					plancheck.ExpectUnknownValue(resourceName, regionConfig0_0.AtMapKey("read_only_specs").AtMapKey("disk_size_gb")), // disk_size_gb is changed so always marked as unknown
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("electable_specs").AtMapKey("disk_size_gb"), knownvalue.Int64Exact(99)),
				},
			},
			{
				ConfigFilename: "main_tls_cipher_config_mode_changed.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectUnknownValue(resourceName, advConfig.AtMapKey("custom_openssl_cipher_config_tls12")),
					plancheck.ExpectKnownValue(resourceName, advConfig.AtMapKey("javascript_enabled"), knownvalue.Bool(true)),
				},
			},
			{
				ConfigFilename: "main_cluster_type_changed.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("config_server_type")),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("config_server_management_mode")),
					plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("bi_connector_config"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"enabled":         knownvalue.Bool(false),
						"read_preference": knownvalue.StringExact("secondary"),
					})),
				},
			},
			{
				ConfigFilename: "main_pin_fcv.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("pinned_fcv").AtMapKey("version")),
				},
			},
		}
	)
	unit.RunPlanCheckTests(t, baseConfig, testCases)
}

func TestPlanChecksClusterTwoRepSpecsMultipleRegions(t *testing.T) {
	var (
		baseConfig   = unit.NewMockPlanChecksConfig(t, &mockConfig, unit.ImportNameTwoRepSpecsMultipleRegions)
		resourceName = baseConfig.ResourceName
		testCases    = []unit.PlanCheckTest{
			{
				ConfigFilename: "main_specs_removed.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("read_only_specs"), specInstanceSizeNodeCount("M10", 2)), // copied from state
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("electable_specs"), specInstanceSizeNodeCount("M10", 5)), // copied from state
					plancheck.ExpectKnownValue(resourceName, regionConfig1_0.AtMapKey("read_only_specs"), specInstanceSizeNodeCount("M20", 1)), // copied from plan in same region_config
					plancheck.ExpectKnownValue(resourceName, regionConfig1_1.AtMapKey("read_only_specs"), specInstanceSizeNodeCount("M20", 2)), // copied from plan in different region_config
					plancheck.ExpectUnknownValue(resourceName, regionConfig1_1.AtMapKey("read_only_specs").AtMapKey("disk_iops")),              // auto_scaling.compute_enabled is used, so disk_iops should be unknown
					plancheck.ExpectUnknownValue(resourceName, regionConfig1_1.AtMapKey("read_only_specs").AtMapKey("disk_size_gb")),           // auto_scaling.disk_enabled is used, so disk_size_gb should be unknown
					plancheck.ExpectUnknownValue(resourceName, regionConfig1_1.AtMapKey("electable_specs")),                                    // node_count = 0, should not be copied
				},
			},
		}
	)
	unit.RunPlanCheckTests(t, baseConfig, testCases)
}

func TestPlanChecksClusterReplicasetNoAutoScaling(t *testing.T) {
	var (
		baseConfig   = unit.NewMockPlanChecksConfig(t, &mockConfig, unit.ImportNameClusterReplicasetNoAutoScaling)
		resourceName = baseConfig.ResourceName
		testCases    = []unit.PlanCheckTest{
			{
				ConfigFilename: "main_region_config_added.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("auto_scaling"), autoScalingKnownValue(false, false, false, "", "")), // using explicit empty values
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("analytics_auto_scaling"), knownvalue.Null()),                        // using explicit null
					plancheck.ExpectKnownValue(resourceName, regionConfig0_1.AtMapKey("auto_scaling"), autoScalingKnownValue(false, false, false, "", "")), // using explicit empty values
					plancheck.ExpectKnownValue(resourceName, regionConfig0_1.AtMapKey("analytics_auto_scaling"), knownvalue.Null()),                        // using explicit null
				},
			},
			{
				ConfigFilename: "main_replication_spec_added.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("auto_scaling"), autoScalingKnownValue(false, false, false, "", "")), // using explicit empty values
					plancheck.ExpectKnownValue(resourceName, regionConfig0_0.AtMapKey("analytics_auto_scaling"), knownvalue.Null()),                        // using explicit null
					plancheck.ExpectKnownValue(resourceName, regionConfig1_0.AtMapKey("auto_scaling"), autoScalingKnownValue(false, false, false, "", "")), // using explicit empty values
					plancheck.ExpectKnownValue(resourceName, regionConfig1_0.AtMapKey("analytics_auto_scaling"), knownvalue.Null()),
				},
			},
		}
	)
	unit.RunPlanCheckTests(t, baseConfig, testCases)
}

package advancedcluster_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

var (
	repSpec0      = tfjsonpath.New("replication_specs").AtSliceIndex(0)
	repSpec1      = tfjsonpath.New("replication_specs").AtSliceIndex(1)
	regionConfig0 = repSpec0.AtMapKey("region_configs").AtSliceIndex(0)
	regionConfig1 = repSpec1.AtMapKey("region_configs").AtSliceIndex(0)
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
				// POC (unsetting O+C sizing/autoscaling attrs): previously a Noop. With the
				// keep-unknown-on-removal plan modifiers, the omitted (computed) disk_iops/disk_size_gb
				// on the kept electable_specs are forced unknown, so this is now an Update. This is the
				// concrete cost of the "can't distinguish removed vs never-set" limitation: fields the
				// user never set in config now show (known after apply) on every plan.
				ConfigFilename: "main_removed_blocks_from_config_no_plan_changes.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectUnknownValue(resourceName, regionConfig0.AtMapKey("electable_specs").AtMapKey("disk_size_gb")),
					plancheck.ExpectUnknownValue(resourceName, regionConfig0.AtMapKey("electable_specs").AtMapKey("disk_iops")),
				},
			},
			{
				// POC: deliberately removing a set autoscaling field (auto_scaling.compute_max_instance_size,
				// "M30" in state) makes the unset a real, sendable change instead of a silent no-op.
				ConfigFilename: "main_removed_optional_attrs.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectUnknownValue(resourceName, regionConfig0.AtMapKey("auto_scaling").AtMapKey("compute_max_instance_size")),
					plancheck.ExpectUnknownValue(resourceName, regionConfig0.AtMapKey("electable_specs").AtMapKey("disk_size_gb")),
				},
			},
			{
				ConfigFilename: "main_node_count_unknown.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("read_only_specs").AtMapKey("node_count"), knownvalue.Int64Exact(2)),
				},
			},
			{
				ConfigFilename: "main_removed_blocks_from_config_and_instance_change.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					// checks regionConfig0
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("read_only_specs"), specInstanceSizeNodeCount("M10", 2)),
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("electable_specs"), specInstanceSizeNodeCount("M10", 5)),
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("auto_scaling"), autoScalingEnabled),
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("analytics_auto_scaling"), autoScalingEnabled),
					plancheck.ExpectUnknownValue(resourceName, regionConfig0.AtMapKey("analytics_specs")), // analytics specs was defined in region_configs.0 but not in region_configs.1

					// checks regionConfig1
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("read_only_specs"), specInstanceSizeNodeCount("M20", 1)),
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("electable_specs"), specInstanceSizeNodeCount("M20", 3)),
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("auto_scaling"), autoScalingEnabled),
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("analytics_auto_scaling"), autoScalingEnabled),
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("analytics_specs"), knownvalue.NotNull()),
				},
			},
		}
	)
	for _, testCase := range testCases {
		t.Run(testCase.ConfigFilename, func(t *testing.T) {
			unit.MockPlanChecksAndRun(t, baseConfig.WithPlanCheckTest(testCase))
		})
	}
}

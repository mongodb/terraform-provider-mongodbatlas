package advancedclustertpf_test

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
	mockConfig    = unit.MockConfigAdvancedClusterTPF
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
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("read_only_specs").AtMapKey("instance_size"), knownvalue.StringExact("M10")),
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("read_only_specs").AtMapKey("instance_size"), knownvalue.StringExact("M20")),
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("auto_scaling"), autoScalingEnabled),
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("analytics_auto_scaling"), autoScalingEnabled),
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("auto_scaling"), autoScalingEnabled),
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("analytics_auto_scaling"), autoScalingEnabled),
					plancheck.ExpectUnknownValue(resourceName, regionConfig0.AtMapKey("analytics_specs")), // analytics specs was defined in region_configs.0 but not in region_configs.1
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("analytics_specs"), knownvalue.NotNull()),
					plancheck.ExpectUnknownValue(resourceName, repSpec0.AtMapKey("id")),
					plancheck.ExpectUnknownValue(resourceName, repSpec1.AtMapKey("id")),
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

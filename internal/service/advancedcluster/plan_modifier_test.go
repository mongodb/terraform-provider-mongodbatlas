package advancedcluster_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
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
				ConfigFilename: "main_removed_blocks_from_config_no_plan_changes.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
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

func TestAdvancedCluster_WarnIgnoredSpecChange(t *testing.T) {
	const rc00 = "replication_specs[0].region_configs[0]"

	testCases := map[string]struct {
		attributeChanges []string
		configs          []advancedcluster.RegionAutoScaling
		expectWarning    bool
	}{
		"warns when compute auto-scaling on and electable instance_size changed": {
			attributeChanges: []string{rc00 + ".electable_specs.instance_size"},
			configs:          []advancedcluster.RegionAutoScaling{{RCPrefix: rc00, ComputeEnabled: true}},
			expectWarning:    true,
		},
		"warns when disk auto-scaling on and electable disk fields changed": {
			attributeChanges: []string{rc00 + ".electable_specs.disk_size_gb", rc00 + ".electable_specs.disk_iops"},
			configs:          []advancedcluster.RegionAutoScaling{{RCPrefix: rc00, DiskGBEnabled: true}},
			expectWarning:    true,
		},
		"warns when compute auto-scaling on and electable disk_size_gb changed": {
			attributeChanges: []string{rc00 + ".electable_specs.disk_size_gb"},
			configs:          []advancedcluster.RegionAutoScaling{{RCPrefix: rc00, ComputeEnabled: true}},
			expectWarning:    true,
		},
		"warns when disk auto-scaling on and electable instance_size changed": {
			attributeChanges: []string{rc00 + ".electable_specs.instance_size"},
			configs:          []advancedcluster.RegionAutoScaling{{RCPrefix: rc00, DiskGBEnabled: true}},
			expectWarning:    true,
		},
		"warns when compute auto-scaling on and read_only instance_size changed": {
			attributeChanges: []string{rc00 + ".read_only_specs.instance_size"},
			configs:          []advancedcluster.RegionAutoScaling{{RCPrefix: rc00, ComputeEnabled: true}},
			expectWarning:    true,
		},
		"warns when analytics compute auto-scaling on and analytics instance_size changed": {
			attributeChanges: []string{rc00 + ".analytics_specs.instance_size"},
			configs:          []advancedcluster.RegionAutoScaling{{RCPrefix: rc00, AnalyticsComputeEnabled: true}},
			expectWarning:    true,
		},
		"no warning when auto-scaling is disabled": {
			attributeChanges: []string{rc00 + ".electable_specs.instance_size"},
			configs:          []advancedcluster.RegionAutoScaling{{RCPrefix: rc00}},
			expectWarning:    false,
		},
		"no warning when auto-scaling is on but no managed spec fields changed": {
			attributeChanges: nil,
			configs:          []advancedcluster.RegionAutoScaling{{RCPrefix: rc00, ComputeEnabled: true}},
			expectWarning:    false,
		},
		"no warning when auto-scaling is on but only node_count changed": {
			attributeChanges: []string{rc00 + ".electable_specs.node_count"},
			configs:          []advancedcluster.RegionAutoScaling{{RCPrefix: rc00, ComputeEnabled: true}},
			expectWarning:    false,
		},
		"no warning when only analytics compute auto-scaling on but electable instance_size changed": {
			attributeChanges: []string{rc00 + ".electable_specs.instance_size"},
			configs:          []advancedcluster.RegionAutoScaling{{RCPrefix: rc00, AnalyticsComputeEnabled: true}},
			expectWarning:    false,
		},
		"no warning when only electable auto-scaling on but analytics instance_size changed": {
			attributeChanges: []string{rc00 + ".analytics_specs.instance_size"},
			configs:          []advancedcluster.RegionAutoScaling{{RCPrefix: rc00, ComputeEnabled: true}},
			expectWarning:    false,
		},
		"no warning when only electable auto-scaling on but analytics disk_size_gb changed": {
			attributeChanges: []string{rc00 + ".analytics_specs.disk_size_gb"},
			configs:          []advancedcluster.RegionAutoScaling{{RCPrefix: rc00, DiskGBEnabled: true}},
			expectWarning:    false,
		},
		"no warning when analytics compute auto-scaling on and analytics disk_size_gb changed": {
			attributeChanges: []string{rc00 + ".analytics_specs.disk_size_gb"},
			configs:          []advancedcluster.RegionAutoScaling{{RCPrefix: rc00, AnalyticsComputeEnabled: true}},
			expectWarning:    false,
		},
		"no warning when analytics instance_size changed without analytics compute auto-scaling": {
			// Confirmed via repro: analytics disk_gb_enabled alone does not cause Atlas to ignore instance_size.
			attributeChanges: []string{rc00 + ".analytics_specs.instance_size"},
			configs:          []advancedcluster.RegionAutoScaling{{RCPrefix: rc00}},
			expectWarning:    false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var diags diag.Diagnostics
			advancedcluster.WarnIgnoredSpecChange(&diags, true, tc.attributeChanges, tc.configs)
			assert.False(t, diags.HasError())
			if tc.expectWarning {
				assert.Equal(t, 1, diags.WarningsCount())
				assert.Contains(t, diags[0].Summary(), "Spec change ignored")
			} else {
				assert.Equal(t, 0, diags.WarningsCount())
			}
		})
	}

	t.Run("no warning when use_effective_fields is false", func(t *testing.T) {
		var diags diag.Diagnostics
		configs := []advancedcluster.RegionAutoScaling{{RCPrefix: rc00, ComputeEnabled: true}}
		advancedcluster.WarnIgnoredSpecChange(&diags, false, []string{rc00 + ".electable_specs.instance_size"}, configs)
		assert.False(t, diags.HasError())
		assert.Equal(t, 0, diags.WarningsCount())
	})
}

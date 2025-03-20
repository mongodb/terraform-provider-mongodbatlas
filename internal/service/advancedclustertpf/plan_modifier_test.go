package advancedclustertpf_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

var (
	repSpec0      = tfjsonpath.New("replication_specs").AtSliceIndex(0)
	repSpec1      = tfjsonpath.New("replication_specs").AtSliceIndex(1)
	regionConfig0 = repSpec0.AtMapKey("region_configs").AtSliceIndex(0)
	regionConfig1 = repSpec1.AtMapKey("region_configs").AtSliceIndex(0)
	mockConfig    = unit.MockHTTPDataConfig{AllowMissingRequests: true, SideEffect: shortenRetries, IsDiffMustSubstrings: []string{"/clusters"}, QueryVars: []string{"providerName"}}
)

func shortenRetries() error {
	advancedclustertpf.RetryMinTimeout = 100 * time.Millisecond
	advancedclustertpf.RetryDelay = 100 * time.Millisecond
	advancedclustertpf.RetryPollInterval = 100 * time.Millisecond
	return nil
}

func autoScalingKnownValue(computeEnabled, diskEnabled, scaleDown bool, minInstanceSize, maxInstanceSize string) knownvalue.Check {
	return knownvalue.ObjectExact(map[string]knownvalue.Check{
		"compute_enabled":            knownvalue.Bool(computeEnabled),
		"disk_gb_enabled":            knownvalue.Bool(diskEnabled),
		"compute_scale_down_enabled": knownvalue.Bool(scaleDown),
		"compute_min_instance_size":  knownvalue.StringExact(minInstanceSize),
		"compute_max_instance_size":  knownvalue.StringExact(maxInstanceSize),
	})
}

func TestMockPlanChecks_ClusterTwoRepSpecsWithAutoScalingAndSpecs(t *testing.T) {
	var (
		baseConfig = unit.NewMockPlanChecksConfig(t, mockConfig, unit.ImportNameClusterTwoRepSpecsWithAutoScalingAndSpecs)
		resourceName = baseConfig.ResourceName
	)
	testCases := map[string][]plancheck.PlanCheck{
		"removed_blocks_from_config_and_instance_change": {
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
			plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("read_only_specs").AtMapKey("instance_size"), knownvalue.StringExact("M10")),
			plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("read_only_specs").AtMapKey("instance_size"), knownvalue.StringExact("M20")),
			plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("auto_scaling"), autoScalingKnownValue(true, true, true, "M10", "M30")),
			plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("analytics_auto_scaling"), autoScalingKnownValue(true, true, true, "M10", "M30")),
			plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("auto_scaling"), autoScalingKnownValue(true, true, true, "M10", "M30")),
			plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("analytics_auto_scaling"), autoScalingKnownValue(true, true, true, "M10", "M30")),
			plancheck.ExpectUnknownValue(resourceName, regionConfig0.AtMapKey("analytics_specs")),
			plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("analytics_specs"), knownvalue.NotNull()),
			plancheck.ExpectUnknownValue(resourceName, repSpec0.AtMapKey("id")),
			plancheck.ExpectUnknownValue(resourceName, repSpec1.AtMapKey("id")),
		},
		"removed_blocks_from_config_no_plan_changes": {
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
		},
	}
	for name, checks := range testCases {
		t.Run(name, func(t *testing.T) {
			unit.MockPlanChecksAndRun(t, baseConfig.WithNameAndChecks(name, checks))
		})
	}
}

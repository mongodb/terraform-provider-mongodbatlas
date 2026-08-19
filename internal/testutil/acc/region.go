package acc

import (
	"os"
	"slices"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
)

// TestRegionsEnvName allows overriding the Atlas regions used by acceptance tests, e.g. to
// avoid OUT_OF_CAPACITY errors in cloud-dev. The value is a comma-separated list of Atlas
// region names (e.g. "US_WEST_2,US_EAST_1,US_EAST_2"). The first region is used by default
// and the rest are fallback candidates for helpers creating clusters directly via the API.
// The first region must have capacity for M10 clusters, and for NVMe-capable tests
// (e.g. M40_NVME) it must also support NVMe instance sizes.
const TestRegionsEnvName = "MONGODB_ATLAS_TEST_REGIONS"

var defaultRegions = []string{constant.UsWest2, constant.UsEast1, constant.UsEast2}

// DefaultRegions returns the Atlas regions used by acceptance test helpers, the first one
// being the default and the rest fallback candidates. It can be overridden with the
// MONGODB_ATLAS_TEST_REGIONS environment variable.
func DefaultRegions() []string {
	regions := strings.Split(os.Getenv(TestRegionsEnvName), ",")
	result := make([]string, 0, len(regions))
	for _, region := range regions {
		if region = strings.TrimSpace(region); region != "" {
			result = append(result, region)
		}
	}
	if len(result) == 0 {
		return slices.Clone(defaultRegions)
	}
	return result
}

// DefaultRegion returns the Atlas region used by default in acceptance test helpers.
func DefaultRegion() string {
	return DefaultRegions()[0]
}

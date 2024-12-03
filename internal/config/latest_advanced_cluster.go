package config

import (
	"os"
	"strconv"
)

const LatestAdvancedClusterEnabledEnvVar = "MONGODB_ATLAS_LATEST_ADVANCED_CLUSTER_ENABLED"
const allowLatestAdvancedClusterEnabled = false // Don't allow in master branch yet, not in const block to allow automatic change

func LatestAdvancedClusterEnabled() bool {
	env, _ := strconv.ParseBool(os.Getenv(LatestAdvancedClusterEnabledEnvVar))
	return allowLatestAdvancedClusterEnabled && env
}

package config

import (
	"os"
	"strconv"
)

const (
	LatestAdvancedClusterEnabledEnvVar = "MONGODB_ATLAS_LATEST_ADVANCED_CLUSTER_ENABLED"
	allowLatestAdvancedClusterEnabled  = true // Don't allow in master branch yet
)

func LatestAdvancedClusterEnabled() bool {
	env, _ := strconv.ParseBool(os.Getenv(LatestAdvancedClusterEnabledEnvVar))
	return allowLatestAdvancedClusterEnabled && env
}

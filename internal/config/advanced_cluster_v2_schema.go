package config

import (
	"os"
	"strconv"
)

const AdvancedClusterV2SchemaEnvVar = "MONGODB_ATLAS_ADVANCED_CLUSTER_V2_SCHEMA"
const allowAdvancedClusterV2Schema = false // Don't allow in master branch yet, not in const block to allow automatic change

func AdvancedClusterV2Schema() bool {
	env, _ := strconv.ParseBool(os.Getenv(AdvancedClusterV2SchemaEnvVar))
	return allowAdvancedClusterV2Schema && env
}

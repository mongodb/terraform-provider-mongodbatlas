package config

import (
	"os"
	"strconv"
)

const AdvancedClusterV2SchemaEnvVar = "MONGODB_ATLAS_ADVANCED_CLUSTER_V2_SCHEMA"
const allowAdvancedClusterV2Schema = false // Don't allow in master branch yet, not in const block to allow automatic change

// Environment variable is read only once to avoid possible changes during runtime
var advancedClusterV2Schema, _ = strconv.ParseBool(os.Getenv(AdvancedClusterV2SchemaEnvVar))

func AdvancedClusterV2Schema() bool {
	return allowAdvancedClusterV2Schema && advancedClusterV2Schema
}

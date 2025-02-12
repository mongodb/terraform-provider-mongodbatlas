package config

import (
	"os"
	"strconv"
)

const PreviewProviderV2AdvancedClusterEnvVar = "MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER"
const allowPreviewProviderV2AdvancedCluster = false // Don't allow in master branch yet, not in const block to allow automatic change

// Environment variable is read only once to avoid possible changes during runtime
var previewProviderV2AdvancedCluster, _ = strconv.ParseBool(os.Getenv(PreviewProviderV2AdvancedClusterEnvVar))

func PreviewProviderV2AdvancedCluster() bool {
	return allowPreviewProviderV2AdvancedCluster && previewProviderV2AdvancedCluster
}

package config

import (
	"os"
	"strconv"
)

const PreviewProviderV2EnvVar = "MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ENABLED"
const allowPreviewProviderV2 = true // Don't allow in master branch yet, not in const block to allow automatic change

// Environment variable is read only once to avoid possible changes during runtime
var previewProviderV2, _ = strconv.ParseBool(os.Getenv(PreviewProviderV2EnvVar))

func PreviewProviderV2() bool {
	return allowPreviewProviderV2 && previewProviderV2
}

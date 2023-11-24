package mongodbatlas_test

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func decodeStateID(stateID string) map[string]string {
	return config.DecodeStateID(stateID)
}

func encodeStateID(values map[string]string) string {
	return config.EncodeStateID(values)
}

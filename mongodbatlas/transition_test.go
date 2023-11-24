package mongodbatlas_test

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errorGetRead = "error reading cloud provider access %s"
)

func decodeStateID(stateID string) map[string]string {
	return config.DecodeStateID(stateID)
}

func encodeStateID(values map[string]string) string {
	return config.EncodeStateID(values)
}

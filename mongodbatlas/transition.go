package mongodbatlas

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

type MongoDBClient = config.MongoDBClient

func encodeStateID(values map[string]string) string {
	return config.EncodeStateID(values)
}

func getEncodedID(stateID, keyPosition string) string {
	return config.GetEncodedID(stateID, keyPosition)
}

func decodeStateID(stateID string) map[string]string {
	return config.DecodeStateID(stateID)
}

package mongodbatlas

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

type MongoDBClient = config.MongoDBClient

func EncodeStateID(values map[string]string) string {
	return config.EncodeStateID(values)
}

func GetEncodedID(stateID, keyPosition string) string {
	return config.GetEncodedID(stateID, keyPosition)
}

func DecodeStateID(stateID string) map[string]string {
	return config.DecodeStateID(stateID)
}

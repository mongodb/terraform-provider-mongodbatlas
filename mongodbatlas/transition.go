package mongodbatlas

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	DeprecationMessageParameterToResource = config.DeprecationMessageParameterToResource
	DeprecationByDateMessageParameter     = config.DeprecationByDateMessageParameter
	DeprecationByDateWithReplacement      = config.DeprecationByDateWithReplacement
	DeprecationByVersionMessageParameter  = config.DeprecationByVersionMessageParameter
	DeprecationMessage                    = config.DeprecationMessage
	AWS                                   = config.AWS
	AZURE                                 = config.AZURE
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

func testAccPreCheckBasic(tb testing.TB) {
	acc.PreCheckBasic(tb)
}

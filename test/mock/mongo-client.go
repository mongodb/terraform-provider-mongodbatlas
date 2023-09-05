package mock

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/config"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func NewMockMongoDBClient() *config.MongoDBClient {
	return &config.MongoDBClient{
		Atlas: &matlas.Client{
			DatabaseUsers: &DatabaseUsersServiceOp{},
		},
	}
}

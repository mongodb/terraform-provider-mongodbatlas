package mongodbatlas

import client "github.com/mongodb-partners/go-client-mongodbatlas/mongodbatlas"

//Config ...
type Config struct {
	PublicKey  string
	PrivateKey string
}

//NewClient ...
func (c *Config) NewClient() interface{} {
	//Initialize the MongoDB Atlas API Client.
	return client.NewClient(nil)
}

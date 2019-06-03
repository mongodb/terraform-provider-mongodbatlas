package mongodbatlas

import "net/http"
import client "github.com/mongodb-partners/go-client-mongodb-atlas/mongodbatlas"
import dac "github.com/xinsnake/go-http-digest-auth-client"

//Config ...
type Config struct {
	PublicKey  string
	PrivateKey string
}

//NewClient ...
func (c *Config) NewClient() interface{} {
	t := dac.NewTransport(c.PublicKey, c.PrivateKey)

	defautlClient := http.DefaultClient
	defautlClient.Transport = &t

	//Initialize the MongoDB Atlas API Client.
	return client.NewClient(defautlClient)
}

package mongodbatlas

import (
	digest "github.com/Sectorbob/mlab-ns2/gae/ns/digest"
	"github.com/hashicorp/terraform/helper/logging"
	matlasClient "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

//Config ...
type Config struct {
	PublicKey  string
	PrivateKey string
}

//NewClient ...
func (c *Config) NewClient() interface{} {
	// setup a transport to handle digest
	transport := digest.NewTransport(c.PublicKey, c.PrivateKey)

	// initialize the client
	client, err := transport.Client()
	if err != nil {
		return err
	}

	client.Transport = logging.NewTransport("MongoDB Atlas", transport)

	//Initialize the MongoDB Atlas API Client.
	return matlasClient.NewClient(client)
}

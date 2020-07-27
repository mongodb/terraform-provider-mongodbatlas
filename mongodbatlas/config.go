package mongodbatlas

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	digest "github.com/mongodb-forks/digest"
	matlasClient "go.mongodb.org/atlas/mongodbatlas"
)

// Config struct ...
type Config struct {
	PublicKey  string
	PrivateKey string
}

// NewClient func...
func (c *Config) NewClient() interface{} {
	// setup a transport to handle digest
	transport := digest.NewTransport(c.PublicKey, c.PrivateKey)

	// initialize the client
	client, err := transport.Client()
	if err != nil {
		return err
	}

	client.Transport = logging.NewTransport("MongoDB Atlas", transport)

	// Initialize the MongoDB Atlas API Client.
	atlasClient := matlasClient.NewClient(client)
	atlasClient.UserAgent = "terraform-provider-mongodbatlas/" + ProviderVersion

	return atlasClient
}

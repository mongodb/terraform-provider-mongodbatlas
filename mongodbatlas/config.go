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
	BaseURL    string
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

	opts := []matlasClient.ClientOpt{matlasClient.SetUserAgent("terraform-provider-mongodbatlas/" + ProviderVersion)}
	if c.BaseURL != "" {
		opts = append(opts, matlasClient.SetBaseURL(c.BaseURL))
	}

	// Initialize the MongoDB Atlas API Client.
	atlasClient, err := matlasClient.New(client, opts...)
	if err != nil {
		return err
	}

	return atlasClient
}

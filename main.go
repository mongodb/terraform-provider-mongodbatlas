package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mongodbatlas.Provider})
}

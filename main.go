package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: mongodbatlas.Provider}

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/mongodb/mongodbatlas", opts)

		if err != nil {
			log.Fatal(err.Error())
		}

		return
	}

	plugin.Serve(opts)
}

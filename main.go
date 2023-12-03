package main

import (
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	var serveOpts []tf6server.ServeOpt

	if debugMode {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}
	err := tf6server.Serve(
		"registry.terraform.io/mongodb/mongodbatlas",
		provider.MuxedProviderFactory(),
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}

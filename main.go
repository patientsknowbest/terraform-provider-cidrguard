package main

import (
	"context"
	"log"

	"github.com/adacasolutions/terraform-provider-cidr-guard/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"
)

func main() {
	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "adacasolutions/cidr-guard",
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}

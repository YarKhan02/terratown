package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), NewProvider, providerserver.ServeOpts{
		Address: "local.providers/local/terratowns",
	})

	if err != nil {
		log.Fatal(err)
	}
}

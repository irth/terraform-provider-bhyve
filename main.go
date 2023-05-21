package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/irth/terraform-provider-bhyve/bhyve"
)

func main() {
	providerserver.Serve(context.Background(), bhyve.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/irth/bhyve",
	})
}

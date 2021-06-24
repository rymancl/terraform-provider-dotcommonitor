package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/rymancl/terraform-provider-dotcommonitor/dotcommonitor"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dotcommonitor.Provider})
}

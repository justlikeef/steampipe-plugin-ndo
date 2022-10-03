package main

import (
	"steampipe-plugin-ndo/ndo"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: ndo.Plugin})
}

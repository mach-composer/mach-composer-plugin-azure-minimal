package main

import (
	"github.com/mach-composer/mach-composer-plugin-sdk/v2/plugin"

	"github.com/mach-composer/mach-composer-plugin-azure-minimal/internal"
)

func main() {
	p := internal.NewAzurePlugin()
	plugin.ServePlugin(p)
}

package internal

import (
	"embed"
	"encoding/json"

	"github.com/mach-composer/mach-composer-plugin-sdk/v2/schema"
)

//go:embed schemas/*
var schemas embed.FS

func getSchema() *schema.ValidationSchema {
	s := schema.ValidationSchema{}
	loadSchemaNode("schemas/site-config.json", &s.SiteConfigSchema)
	loadSchemaNode("schemas/remote-state.json", &s.RemoteStateSchema)
	loadSchemaNode("schemas/global-config.json", &s.GlobalConfigSchema)

	// These are the same (intended)
	loadSchemaNode("schemas/component-config.json", &s.ComponentConfigSchema)
	loadSchemaNode("schemas/component-config.json", &s.SiteComponentConfigSchema)

	return &s
}

func loadSchemaNode(filename string, dst any) {
	body, err := schemas.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, dst); err != nil {
		panic(err)
	}
}

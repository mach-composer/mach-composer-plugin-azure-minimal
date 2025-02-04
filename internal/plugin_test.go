package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlugin_RenderTerraformComponent_OK(t *testing.T) {
	var component = &ComponentConfig{
		ServicePlan: "plan",
	}

	p := &Plugin{
		globalConfig: &GlobalConfig{
			SubscriptionID: "subscription",
			ResourceGroup:  "name",
			ResourcePrefix: "prefix",
		},
		siteConfigs: map[string]SiteConfig{
			"site": {
				Components: map[string]SiteComponentConfig{
					"component": {
						Component: component,
					},
				},
			},
		},
		componentConfigs: map[string]ComponentConfig{
			"component": *component,
		},
	}
	s, err := p.RenderTerraformComponent("site", "component")

	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, []string{"azurerm = azurerm"}, s.Providers)
	assert.Equal(t, "\n\tazure = {\n\t\tresource_group_name = \"name\"\n\t\tresource_prefix = \"prefix\"\n\t}\n\t", s.Variables)
}

func TestPlugin_RenderTerraformComponent_NoVariables(t *testing.T) {
	var component = &ComponentConfig{
		ServicePlan: "plan",
	}

	p := &Plugin{
		globalConfig: &GlobalConfig{
			SubscriptionID: "subscription",
		},
		siteConfigs: map[string]SiteConfig{
			"site": {
				Components: map[string]SiteComponentConfig{
					"component": {
						Component: component,
					},
				},
			},
		},
		componentConfigs: map[string]ComponentConfig{
			"component": *component,
		},
	}
	_, err := p.RenderTerraformComponent("site", "component")
	assert.Error(t, err)
}

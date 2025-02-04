package internal

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/hashicorp/go-hclog"
	"github.com/mach-composer/mach-composer-plugin-helpers/helpers"
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/mitchellh/mapstructure"
)

func NewAzurePlugin() schema.MachComposerPlugin {
	state := &Plugin{
		provider:         "3.42.0",
		siteConfigs:      map[string]SiteConfig{},
		componentConfigs: map[string]ComponentConfig{},
	}

	return plugin.NewPlugin(&schema.PluginSchema{
		Identifier: "azure",

		Configure: state.Configure,
		IsEnabled: state.IsEnabled,

		// Schema
		GetValidationSchema: state.GetValidationSchema,

		// Config
		SetRemoteStateBackend:  state.SetRemoteStateBackend,
		SetGlobalConfig:        state.SetGlobalConfig,
		SetSiteConfig:          state.SetSiteConfig,
		SetComponentConfig:     state.SetComponentConfig,
		SetSiteComponentConfig: state.SetSiteComponentConfig,

		// Renders
		RenderTerraformStateBackend: state.TerraformRenderStateBackend,
		RenderTerraformProviders:    state.TerraformRenderProviders,
		RenderTerraformResources:    state.TerraformRenderResources,
		RenderTerraformComponent:    state.RenderTerraformComponent,
	})
}

type Plugin struct {
	environment      string
	provider         string
	remoteState      *AzureTFState
	globalConfig     *GlobalConfig
	siteConfigs      map[string]SiteConfig
	componentConfigs map[string]ComponentConfig
}

func (p *Plugin) Configure(environment string, provider string) error {
	p.environment = environment
	if provider != "" {
		p.provider = provider
	}
	return nil
}

func (p *Plugin) IsEnabled() bool {
	return len(p.siteConfigs) > 0
}

func (p *Plugin) SetRemoteStateBackend(data map[string]any) error {
	state := &AzureTFState{}
	if err := mapstructure.Decode(data, state); err != nil {
		return err
	}
	if err := defaults.Set(state); err != nil {
		return err
	}
	p.remoteState = state
	return nil
}

func (p *Plugin) GetValidationSchema() (*schema.ValidationSchema, error) {
	result := getSchema()
	return result, nil
}

func (p *Plugin) SetGlobalConfig(data map[string]any) error {
	if err := mapstructure.Decode(data, &p.globalConfig); err != nil {
		return err
	}

	if p.globalConfig.ResourceTags != nil {
		hclog.Default().Warn(
			fmt.Sprintf("Using resource tags is deprecated. These should be inferred from the site configuration. The field will be removed in a future release."),
		)
	}

	return nil
}

func (p *Plugin) SetSiteConfig(site string, data map[string]any) error {
	if p.globalConfig == nil {
		return fmt.Errorf("a global azure config is required for setting per-site configuration")
	}

	cfg := SiteConfig{
		Components: make(map[string]SiteComponentConfig),
	}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	cfg.merge(p.globalConfig)

	p.siteConfigs[site] = cfg
	return nil
}

func (p *Plugin) SetSiteComponentConfig(site, component string, data map[string]any) error {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		return nil
	}

	c := SiteComponentConfig{
		Name: component,
	}
	if err := mapstructure.Decode(data, &c); err != nil {
		return err
	}

	cfg.Components[component] = c
	p.siteConfigs[site] = cfg
	return nil
}

func (p *Plugin) SetComponentConfig(component string, data map[string]any) error {
	cfg, ok := p.componentConfigs[component]
	if !ok {
		cfg = ComponentConfig{}
	}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	cfg.Name = component
	p.componentConfigs[component] = cfg
	return nil
}

func (p *Plugin) TerraformRenderStateBackend(site string) (string, error) {
	if p.remoteState == nil {
		return "", nil
	}

	templateContext := struct {
		State *AzureTFState
		Site  string
		Key   string
	}{
		State: p.remoteState,
		Site:  site,
		Key:   p.remoteState.Key(site),
	}

	template := `
	backend "azurerm" {
	  resource_group_name  = "{{ .State.ResourceGroup }}"
	  storage_account_name = "{{ .State.StorageAccount }}"
	  container_name       = "{{ .State.ContainerName }}"
	  key                  = "{{ .Key }}"
	}
	`
	return helpers.RenderGoTemplate(template, templateContext)
}

func (p *Plugin) TerraformRenderProviders(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	result := fmt.Sprintf(`
		azurerm = {
			version = "%s"
		}`, helpers.VersionConstraint(p.provider))
	return result, nil
}

func (p *Plugin) TerraformRenderResources(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	var tags = map[string]string{
		"SiteName":    site,
		"Environment": p.environment,
	}
	if p.globalConfig.ResourceTags != nil {
		for k, v := range p.globalConfig.ResourceTags {
			tags[k] = v
		}
	}

	templateContext := struct {
		SubscriptionID string
		Tags           map[string]string
	}{
		SubscriptionID: cfg.SubscriptionID,
		Tags:           tags,
	}

	template := `
		provider "azurerm" {
			subscription_id            = "{{ .SubscriptionID }}"
			skip_provider_registration = true

			features {
				resource_group {
					prevent_deletion_if_contains_resources = true
				}
				key_vault {
					purge_soft_deleted_keys_on_destroy = true
					recover_soft_deleted_keys          = true
				}
			}
		}

		locals {
			{{ renderProperty "tags" .Tags }}
		}
	`
	return helpers.RenderGoTemplate(template, templateContext)
}

func (p *Plugin) RenderTerraformComponent(site string, component string) (*schema.ComponentSchema, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return &schema.ComponentSchema{}, nil
	}

	siteComponent, ok := cfg.Components[component]
	if !ok {
		return nil, fmt.Errorf("missing config for component")
	}
	siteComponent.Component = p.getComponentConfig(component)

	result := &schema.ComponentSchema{
		Providers: []string{"azurerm = azurerm"},
	}

	value, err := terraformRenderComponentVars(cfg, &siteComponent)
	if err != nil {
		return nil, err
	}
	result.Variables = value

	result.DependsOn = []string{}
	return result, nil
}

func (p *Plugin) getSiteConfig(site string) *SiteConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		return nil
	}
	cfg.merge(p.globalConfig)
	return &cfg
}

func (p *Plugin) getComponentConfig(name string) *ComponentConfig {
	componentConfig, ok := p.componentConfigs[name]
	if !ok {
		componentConfig = ComponentConfig{}
	}
	return &componentConfig
}

func terraformRenderComponentVars(cfg *SiteConfig, _ *SiteComponentConfig) (string, error) {
	templateContext := struct {
		Config *SiteConfig
	}{
		Config: cfg,
	}

	template := `
	azure = {
		{{ renderProperty "resource_group_name" .Config.ResourceGroup }}
		{{ renderProperty "resource_prefix" .Config.ResourcePrefix }}
	}
	`
	return helpers.RenderGoTemplate(template, templateContext)
}

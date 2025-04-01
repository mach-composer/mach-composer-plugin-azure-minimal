package internal

import (
	"fmt"
)

// AzureTFState Azure storage account state backend configuration.
type AzureTFState struct {
	ResourceGroup  string `mapstructure:"resource_group"`
	StorageAccount string `mapstructure:"storage_account"`
	ContainerName  string `mapstructure:"container_name"`
	StateFolder    string `mapstructure:"state_folder"`
}

func (a AzureTFState) Key(site string) string {
	if a.StateFolder == "" {
		return site
	}
	return fmt.Sprintf("%s/%s", a.StateFolder, site)
}

type GlobalConfig struct {
	TenantID       string                            `mapstructure:"tenant_id"`
	SubscriptionID string                            `mapstructure:"subscription_id"`
	ResourceGroup  string                            `mapstructure:"resource_group"`
	ResourcePrefix string                            `mapstructure:"resource_prefix"`
	ResourceTags   map[string]string                 `mapstructure:"resource_tags"`
	Features       map[string]map[string]interface{} `mapstructure:"features"`
}

type SiteConfig struct {
	ResourceGroup  string                            `mapstructure:"resource_group"`
	ResourcePrefix string                            `mapstructure:"resource_prefix"`
	SubscriptionID string                            `mapstructure:"subscription_id"`
	Features       map[string]map[string]interface{} `mapstructure:"features"`

	Components map[string]SiteComponentConfig
}

func (s *SiteConfig) merge(g *GlobalConfig) {
	if s.ResourcePrefix == "" {
		s.ResourcePrefix = g.ResourcePrefix
	}
	if s.SubscriptionID == "" {
		s.SubscriptionID = g.SubscriptionID
	}
	if s.ResourceGroup == "" {
		s.ResourceGroup = g.ResourceGroup
	}
	if s.Features == nil {
		s.Features = g.Features
	}
}

type ComponentConfig struct {
	Endpoints map[string]string `mapstructure:"-"`

	Name        string `mapstructure:"-"`
	ServicePlan string `mapstructure:"service_plan"`
	ShortName   string `mapstructure:"short_name"`
}

type SiteComponentConfig struct {
	ServicePlan string `mapstructure:"service_plan"`

	Name      string           `mapstructure:"-"`
	Component *ComponentConfig `mapstructure:"-"`
}

type SiteComponent struct {
	InternalName string
	ExternalName string
	Component    *ComponentConfig
}

package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSiteConfig_merge_WithEmptyFields(t *testing.T) {
	g := &GlobalConfig{
		ResourcePrefix: "globalPrefix",
		SubscriptionID: "globalSubID",
		ResourceGroup:  "globalGroup",
	}
	s := &SiteConfig{}
	s.merge(g)
	assert.Equal(t, "globalPrefix", s.ResourcePrefix)
	assert.Equal(t, "globalSubID", s.SubscriptionID)
	assert.Equal(t, "globalGroup", s.ResourceGroup)
}

func TestSiteConfig_merge_WithNonEmptyFields(t *testing.T) {
	g := &GlobalConfig{
		ResourcePrefix: "globalPrefix",
		SubscriptionID: "globalSubID",
		ResourceGroup:  "globalGroup",
	}
	s := &SiteConfig{
		ResourcePrefix: "sitePrefix",
		SubscriptionID: "siteSubID",
		ResourceGroup:  "siteGroup",
	}
	s.merge(g)
	assert.Equal(t, "sitePrefix", s.ResourcePrefix)
	assert.Equal(t, "siteSubID", s.SubscriptionID)
	assert.Equal(t, "siteGroup", s.ResourceGroup)
}

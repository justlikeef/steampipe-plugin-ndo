package ndo

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type NDOConfig struct {
	ClusterURI         *string `cty:"cluster_uri"`
	AllowUnverifiedSSL *bool   `cty:"allow_unverified_ssl"`
	User               *string `cty:"user"`
	Password           *string `cty:"password"`
	LoginDomain        *string `cty:"login_domain"`
	Platform           *string `cty:"platform"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"cluster_uri": {
		Type: schema.TypeString,
	},
	"allow_unverified_ssl": {
		Type: schema.TypeBool,
	},
	"user": {
		Type: schema.TypeString,
	},
	"password": {
		Type: schema.TypeString,
	},
	"login_domain": {
		Type: schema.TypeString,
	},
	"platform": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &NDOConfig{}
}

func GetConfig(connection *plugin.Connection) NDOConfig {
	if connection == nil || connection.Config == nil {
		return NDOConfig{}
	}

	config, _ := connection.Config.(NDOConfig)
	return config
}

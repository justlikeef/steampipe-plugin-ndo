/*
Package aws implements a steampipe plugin for aci.

This plugin provides data that Steampipe uses to present foreign
tables that represent Cisco ACI resources.
*/
package ndo

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

const pluginName = "steampipe-plugin-ndo"

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-ndo",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		TableMap: map[string]*plugin.Table{
			"ndo_epg_static_port":         tableNDOEpgStaticPort(),
			"ndo_schema":                  tableNDOSchema(),
			"ndo_schema_template":         tableNDOSchemaTemplate(),
			"ndo_schema_template_anp":     tableNDOSchemaTemplateAnp(),
			"ndo_schema_template_vrf":     tableNDOSchemaTemplateVrf(),
			"ndo_schema_template_bd":      tableNDOSchemaTemplateBd(),
			"ndo_schema_template_anp_epg": tableNDOSchemaTemplateAnpEpg(),
		},
	}
	return p
}

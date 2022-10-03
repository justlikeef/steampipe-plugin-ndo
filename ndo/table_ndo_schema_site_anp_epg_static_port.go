package ndo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"steampipe-plugin-ndo/client"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type SchemaSiteAnpEpgStaticPort struct {
	Id                  string
	SchemaId            string
	SiteId              string
	TemplateName        string
	AnpName             string
	EpgName             string
	PathType            string
	Pod                 string
	Leaf                string
	Path                string
	Mode                string
	DeploymentImmediacy string
	Vlan                string
	MicroSegVlan        string
	Fex                 string
}

func tableNDOEpgStaticPort() *plugin.Table {
	return &plugin.Table{
		Name:        "ndo_schema_site_anp_epg_static_port",
		Description: "NDO Schema-Site-ANP-EPG Static Ports",
		List: &plugin.ListConfig{
			Hydrate: listSchemaSiteAnpEpgStaticPorts,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "id",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "(Required) Unique ID of this object within Nexus Dashboard Orchestrator (NDO)",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "schema_id",
				Description: "(Required) SchemaID under which you want to deploy Static Port.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SchemaId"),
			},
			{
				Name:        "site_id",
				Description: "(Required) SiteID under which you want to deploy Static Port.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SiteId"),
			},
			{
				Name:        "template_name",
				Description: "(Required) Template name under which you want to deploy Static Port.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "anp_name",
				Description: "(Required) ANP name under which you want to deploy Static Port.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "epg_name",
				Description: "(Required) EPG name under which you want to deploy Static Port.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "path_type",
				Description: "(Required) The type of the static port. Allowed values are port, vpc and dpc.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pod",
				Description: "(Required) The pod of the static port.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "leaf",
				Description: "(Required) The leaf of the static port.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "path",
				Description: "(Required) The path of the static port.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mode",
				Description: "(Required) The mode of the static port. Allowed values are native, regular and untagged.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deployment_immediacy",
				Description: "(Required) The deployment immediacy of the static port. Allowed values are immediate and lazy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vlan",
				Description: "(Required) The port encap VLAN id of the static port.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "micro_seg_vlan",
				Description: "(Optional) The microsegmentation VLAN id of the static port.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "fex",
				Description: "(Optional) Fex-id to be used. This parameter will work only with the path_type as port.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

//// LIST FUNCTION
func listSchemaSiteAnpEpgStaticPorts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	ndoclient, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to ND: %v", err)
	}

	log.Printf("[DEBUG] Calling API: list-identity")
	dnUrl := "/api/v1/schemas/list-identity"
	identityList, err := ndoclient.ServiceManager.GetViaURL(dnUrl)

	if err != nil {
		return nil, fmt.Errorf("Error getting Identity List: %v\nURL: %v", err, dnUrl)
	}

	log.Printf("[TRACE] Indentity List: %v", identityList)
	schemaobjlist, err := identityList.S("schemas").Children()
	if err != nil {
		return nil, fmt.Errorf("Error getting Schema List: %v", err)
	}

	for _, curschema := range schemaobjlist {
		dnUrl := "/api/v1/schemas/" + client.StripQuotes(curschema.S("id").String())
		schemaDetails, err := ndoclient.ServiceManager.GetViaURL(dnUrl)

		sitelist, err := schemaDetails.S("sites").Children()
		if err != nil {
			return nil, fmt.Errorf("Error getting site list: %v", err)
		}

		for _, cursite := range sitelist {
			anpobjlist, err := cursite.S("anps").Children()
			if err != nil {
				return nil, fmt.Errorf("Error getting anp list: %v", err)
			}

			for _, curanp := range anpobjlist {
				epgobjlist, err := curanp.S("epgs").Children()
				if err != nil {
					return nil, fmt.Errorf("Error getting epg list: %v", err)
				}

				for _, curepg := range epgobjlist {
					staticportlist, err := curepg.S("staticPorts").Children()
					if err != nil {
						return nil, fmt.Errorf("Error getting port list: %v", err)
					}

					for _, curport := range staticportlist {
						portobj := &SchemaSiteAnpEpgStaticPort{}
						portobj.SchemaId = client.StripQuotes(curschema.S("id").String())
						portobj.SiteId = client.StripQuotes(cursite.S("siteId").String())
						portobj.TemplateName = client.StripQuotes(cursite.S("templateName").String())
						epgPathInfo := strings.Split(client.StripQuotes(curepg.S("epgRef").String()), "/")
						portobj.AnpName = epgPathInfo[6]
						portobj.EpgName = epgPathInfo[8]
						portobj.PathType = client.StripQuotes(curport.S("type").String())
						portobj.DeploymentImmediacy = client.StripQuotes(curport.S("deploymentImmediacy").String())
						portPathInfo := strings.Split(client.StripQuotes(curport.S("path").String()), "/")
						portobj.Pod = portPathInfo[1]
						portobj.Leaf = portPathInfo[2][10:len(portPathInfo[2])]
						log.Printf("[TRACE] path: %s: length: %v len-1: %v value: %v", portPathInfo[3], len(portPathInfo[3]), len(portPathInfo[3])-1, portPathInfo[3][8:len(portPathInfo[3])-1])
						portobj.Path = portPathInfo[3]
						portobj.Path = portobj.Path[8 : len(portPathInfo[3])-1]
						portobj.Vlan = client.StripQuotes(curport.S("portEncapVlan").String())
						portobj.Mode = client.StripQuotes(curport.S("mode").String())
						portobj.Id = portobj.SchemaId + "/site/" + portobj.SiteId + "/template/" + portobj.TemplateName + "/anp/" + portobj.AnpName + "/epg/" + portobj.EpgName + "/staticPortPod/" + portobj.Pod + "/staticPortLeaf/" + portobj.Leaf + "/pathType/" + portobj.PathType + "/fex/" + portobj.Leaf + "/path/" + portobj.Path
						log.Printf("[TRACE] Record object: %s ", portobj)
						d.StreamListItem(ctx, portobj)
					}
				}
			}
		}
	}

	return nil, nil
}

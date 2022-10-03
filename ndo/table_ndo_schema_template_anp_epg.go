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

type SchemaTemplateAnpEpg struct {
	Id                       string
	SchemaId                 string
	TemplateName             string
	AnpName                  string
	Name                     string
	BdName                   string
	BdSchemaId               string
	BdTemplateName           string
	VrfName                  string
	VrfSchemaId              string
	VrfTemplateName          string
	DisplayName              string
	UsegEpg                  string
	IntraEpg                 string
	IntersiteMulticastSource string
	ProxyArp                 string
	PreferredGroup           string
}

func tableNDOSchemaTemplateAnpEpg() *plugin.Table {
	return &plugin.Table{
		Name:        "ndo_schema_template_anp_epg",
		Description: "NDO Schema-Template-Anp-Epg",
		List: &plugin.ListConfig{
			Hydrate: listSchemaTemplateAnpEpg,
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
				Description: "(Required) SchemaID under which you want to deploy Anp Epg.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SchemaId"),
			},
			{
				Name:        "template_name",
				Description: "(Required) Template where Anp Epg to be created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "anp_name",
				Description: "(Required) Name of Application Network Profile.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "(Required) Name of Endpoint Group to manage.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bd_name",
				Description: "(Optional) Name of Bridge Domain. It is required when using on-premise sites.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bd_schema_id",
				Description: "(Opional) The schemaID that defines the referenced BD.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BdSchemaId"),
			},
			{
				Name:        "bd_template_name",
				Description: "(Optional) The template that defines the referenced BD.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vrf_name",
				Description: "(Optional) Name of Vrf. It is required when using cloud sites.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vrf_schema_id",
				Description: "(Optional) The schemaID that defines the referenced VRF.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VrfSchemaId"),
			},
			{
				Name:        "vrf_template_name",
				Description: "(Optional) The template that defines the referenced VRF.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "(Optional) The name as displayed on the MSO web interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "useg_epg",
				Description: "(Optional) Boolean flag to enable or disable whether this is a USEG EPG. Default value is set to false.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "intra_epg",
				Description: "(Optional) Whether intra EPG isolation is enforced. choices: [ enforced, unenforced ]",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "intersite_multicast_source",
				Description: "(Optional) Whether intersite multicast source is enabled. Default to false.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "proxy_arp",
				Description: "(Optional) Whether to enable Proxy ARP or not. (For Forwarding control) Default to false.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "preferred_group",
				Description: "(Optional) Boolean flag to enable or disable whether this EPG is added to preferred group. Default value is set to false.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

//// LIST FUNCTION
func listSchemaTemplateAnpEpg(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

		teamplatelist, err := schemaDetails.S("templates").Children()
		if err != nil {
			return nil, fmt.Errorf("Error getting site list: %v", err)
		}

		for _, curtemp := range teamplatelist {
			anpobjlist, err := curtemp.S("anps").Children()
			if err != nil {
				return nil, fmt.Errorf("Error getting anp list: %v", err)
			}

			for _, curanp := range anpobjlist {
				epgobjlist, err := curanp.S("epgs").Children()
				if err != nil {
					return nil, fmt.Errorf("Error getting epg list: %v", err)
				}

				for _, curepg := range epgobjlist {
					epgobj := &SchemaTemplateAnpEpg{}
					epgobj.SchemaId = client.StripQuotes(curschema.S("id").String())
					epgobj.TemplateName = client.StripQuotes(curtemp.S("name").String())
					epgobj.AnpName = client.StripQuotes(curanp.S("name").String())
					epgobj.Name = client.StripQuotes(curepg.S("name").String())
					if client.StripQuotes(curepg.S("bdRef").String()) != "" {
						bdInfo := strings.Split(client.StripQuotes(curepg.S("bdRef").String()), "/")
						epgobj.BdName = bdInfo[6]
						epgobj.BdSchemaId = bdInfo[2]
						epgobj.BdTemplateName = bdInfo[4]
					}
					if client.StripQuotes(curepg.S("vrfRef").String()) != "" {
						vrfInfo := strings.Split(client.StripQuotes(curepg.S("vrfRef").String()), "/")
						epgobj.VrfName = vrfInfo[5]
						epgobj.VrfSchemaId = vrfInfo[1]
						epgobj.VrfTemplateName = vrfInfo[3]
					}
					epgobj.DisplayName = client.StripQuotes(curepg.S("displayName").String())
					epgobj.UsegEpg = client.StripQuotes(curepg.S("uSegEpg").String())
					epgobj.IntraEpg = client.StripQuotes(curepg.S("intraEpg").String())
					epgobj.IntersiteMulticastSource = client.StripQuotes(curepg.S("mCastSource").String())
					epgobj.ProxyArp = client.StripQuotes(curepg.S("proxyArp").String())
					epgobj.PreferredGroup = client.StripQuotes(curepg.S("preferredGroup").String())
					epgobj.Id = epgobj.SchemaId + "/template/" + epgobj.TemplateName + "/anp/" + epgobj.AnpName + "/epg/" + epgobj.Name
					log.Printf("[TRACE] Record object: %s ", epgobj)
					d.StreamListItem(ctx, epgobj)
				}
			}
		}
	}
	return nil, nil
}

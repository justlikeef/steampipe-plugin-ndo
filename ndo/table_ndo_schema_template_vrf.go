package ndo

import (
	"context"
	"fmt"
	"log"

	"steampipe-plugin-ndo/client"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type SchemaTemplateVrf struct {
	Id              string
	SchemaId        string
	Name            string
	Template        string
	DisplayName     string
	Layer3Multicast string
	Vzany           string
}

func tableNDOSchemaTemplateVrf() *plugin.Table {
	return &plugin.Table{
		Name:        "ndo_schema_template_vrf",
		Description: "NDO Schema-Template-Vrf",
		List: &plugin.ListConfig{
			Hydrate: listSchemaTemplateVrf,
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
				Description: "(Required) The schema-id where vrf is associated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SchemaId"),
			},
			{
				Name:        "name",
				Description: "(Required) name of the vrf.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "template",
				Description: "(Required) template associated with the vrf.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "(Required) The name as displayed on the MSO web interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "layer3_multicast",
				Description: "(Optional) Whether to enable L3 multicast.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vzany",
				Description: "(Optional) Whether to enable vzany.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

//// LIST FUNCTION
func listSchemaTemplateVrf(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

		templatelist, err := schemaDetails.S("templates").Children()
		if err != nil {
			return nil, fmt.Errorf("Error getting template list: %v", err)
		}

		for _, curtemp := range templatelist {
			vrfobjlist, err := curtemp.S("vrfs").Children()
			if err != nil {
				return nil, fmt.Errorf("Error getting vrf list: %v", err)
			}

			for _, curvrf := range vrfobjlist {
				vrfobj := &SchemaTemplateVrf{}
				vrfobj.SchemaId = client.StripQuotes(curschema.S("id").String())
				vrfobj.Name = client.StripQuotes(curvrf.S("name").String())
				vrfobj.Template = client.StripQuotes(curtemp.S("name").String())
				vrfobj.DisplayName = client.StripQuotes(curvrf.S("displayName").String())
				vrfobj.Layer3Multicast = client.StripQuotes(curvrf.S("l3MCast").String())
				vrfobj.Vzany = client.StripQuotes(curvrf.S("vzAnyEnabled").String())
				vrfobj.Id = vrfobj.SchemaId + "/template/" + vrfobj.Template + "/vrf/" + vrfobj.Name
				log.Printf("[TRACE] Record object: %s ", vrfobj)
				d.StreamListItem(ctx, vrfobj)
			}
		}
	}

	return nil, nil
}

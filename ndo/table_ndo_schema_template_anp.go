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

type SchemaTemplateAnp struct {
	Id          string
	SchemaId    string
	Name        string
	Template    string
	DisplayName string
}

func tableNDOSchemaTemplateAnp() *plugin.Table {
	return &plugin.Table{
		Name:        "ndo_schema_template_anp",
		Description: "NDO Schema-Template-Anp",
		List: &plugin.ListConfig{
			Hydrate: listSchemaTemplateAnp,
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
				Description: "(Required) The schema-id where anp is associated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SchemaId"),
			},
			{
				Name:        "name",
				Description: "(Required) name of the anp.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "template",
				Description: "(Required) template associated with the anp.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "(Required) The name as displayed on the MSO web interface.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

//// LIST FUNCTION
func listSchemaTemplateAnp(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
			anpobjlist, err := curtemp.S("anps").Children()
			if err != nil {
				return nil, fmt.Errorf("Error getting anp list: %v", err)
			}

			for _, curanp := range anpobjlist {
				anpobj := &SchemaTemplateAnp{}
				anpobj.SchemaId = client.StripQuotes(curschema.S("id").String())
				anpobj.Name = client.StripQuotes(curanp.S("name").String())
				anpobj.Template = client.StripQuotes(curtemp.S("name").String())
				anpobj.DisplayName = client.StripQuotes(curanp.S("displayName").String())
				anpobj.Id = anpobj.SchemaId + "/template/" + anpobj.Template + "/anp/" + anpobj.Name
				log.Printf("[TRACE] Record object: %s ", anpobj)
				d.StreamListItem(ctx, anpobj)
			}
		}
	}

	return nil, nil
}

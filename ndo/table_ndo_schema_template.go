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

type SchemaTemplate struct {
	Id          string
	SchemaId    string
	TenantId    string
	Name        string
	DisplayName string
}

func tableNDOSchemaTemplate() *plugin.Table {
	return &plugin.Table{
		Name:        "ndo_schema_template",
		Description: "NDO Schema Template",
		List: &plugin.ListConfig{
			Hydrate: listSchemaTemplate,
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
				Description: "(Required) Name of the schema.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SchemaId"),
			},
			{
				Name:        "tenant_id",
				Description: "(Required) Tenant-id to associate.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TenantId"),
			},
			{
				Name:        "name",
				Description: "(Required) Name of the template.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "(Required) Display name of the Template to be deployed on the site.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

func listSchemaTemplate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	log.Printf("[DEBUG] Getting client")
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
	identityListChildren, err := identityList.S("schemas").Children()
	if err != nil {
		return nil, fmt.Errorf("Error getting Schema List: %v", err)
	}

	for _, curschema := range identityListChildren {
		log.Printf("[TRACE] Processing Schema: %v", curschema)
		templateList, err := curschema.S("templates").Children()
		if err != nil {
			return nil, fmt.Errorf("Error getting Template List: %v", err)
		}

		for _, curtemplate := range templateList {
			log.Printf("[TRACE] Processing Template: %v", curtemplate)
			log.Printf("[TRACE] inside schema: %v", curschema)

			schemaTemplateObj := SchemaTemplate{}
			schemaTemplateObj.Id = client.StripQuotes(curschema.S("id").String()) + "/template/" + client.StripQuotes(curtemplate.S("name").String())
			schemaTemplateObj.SchemaId = client.StripQuotes(curschema.S("id").String())
			schemaTemplateObj.TenantId = client.StripQuotes(curtemplate.S("tenantId").String())
			schemaTemplateObj.Name = client.StripQuotes(curtemplate.S("name").String())
			schemaTemplateObj.DisplayName = client.StripQuotes(curtemplate.S("displayName").String())

			log.Printf("[TRACE] Built Schema-Template Object: %v", schemaTemplateObj)

			d.StreamListItem(ctx, schemaTemplateObj)
		}
	}

	return nil, nil
}

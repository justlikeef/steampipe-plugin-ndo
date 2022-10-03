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

type SchemaTemplateBd struct {
	Id                                string
	SchemaId                          string
	TemplateName                      string
	Name                              string
	DisplayName                       string
	VrfName                           string
	VrfSchemaId                       string
	VrfTemplateName                   string
	Layer2UnknownUnicast              string
	IntersiteBumTraffic               string
	OptimizeWanBandwidth              string
	Layer2Stretch                     string
	Layer3Multicast                   string
	ArpFlooding                       string
	VirtualMacAddress                 string
	UnicastRouting                    string
	Ipv6UnknownMulticastFlooding      string
	MultiDestinationFlooding          string
	UnknownMulticastFlooding          string
	DhcpPolicy                        string
	DhcpPolicyName                    string
	DhcpPolicyVersion                 string
	DhcpPolicyDhcpOptionPolicyName    string
	DhcpPolicyDhcpOptionPolicyVersion string
}

func tableNDOSchemaTemplateBd() *plugin.Table {
	return &plugin.Table{
		Name:        "ndo_schema_template_bd",
		Description: "NDO Schema-Template-Bd",
		List: &plugin.ListConfig{
			Hydrate: listSchemaTemplateBd,
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
				Description: "(Required) SchemaID under which you want to deploy Bridge Domain.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SchemaId"),
			},
			{
				Name:        "template_name",
				Description: "(Required) Template where Bridge Domain to be created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "(Required) Name of Bridge Domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "(Required) Display Name of the Bridge Domain on the MSO UI.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vrf_name",
				Description: "(Required) Name of VRF to attach with Bridge Domain. VRF must exist.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vrf_schema_id",
				Description: "(Optional) SchemaID of VRF. schema_id of Bridge Domain will be used if not provided. Should use this parameter when VRF is in different schema than BD.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VrfSchemaId"),
			},
			{
				Name:        "vrf_template_name",
				Description: "(Optional) Template Name of VRF. template_name of Bridge Domain will be used if not provided. Should use this parameter when VRF is in different schema than BD.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "layer2_unknown_unicast",
				Description: "(Optional) Type of layer 2 unknown unicast. Allowed values are flood and proxy. Default to flood.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "intersite_bum_traffic",
				Description: "(Optional) Boolean Flag to enable or disable intersite bum traffic. Default to false.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "optimize_wan_bandwidth",
				Description: "(Optional) Boolean flag to enable or disable the wan bandwidth optimization. Default to false.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "layer2_stretch",
				Description: "(Optional) Boolean flag to enable or disable the layer-2 stretch. Default to false. Should enable this flag if you want to create subnets under this Bridge Domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "layer3_multicast",
				Description: "(Optional) Boolean flag to enable or disable layer 3 multicast traffic. Default to false.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arp_flooding",
				Description: "(Optional) ARP Flooding status. Default to false.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "virtual_mac_address",
				Description: "(Optional) Virtual MAC Address.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "unicast_routing",
				Description: "(Optional) Unicast Routing status. Default to false.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ipv6_unknown_multicast_flooding",
				Description: "(Optional) IPv6 Unknown Multicast Flooding behavior. Allowed values are flood and optimized_flooding. Default to flood.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "multi_destination_flooding",
				Description: "(Optional) Multi-destination flooding behavior. Allowed values are flood_in_bd, drop and flood_in_encap. Default to flood_in_bd.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "unknown_multicast_flooding",
				Description: "(Optional) Unknown Multicast Flooding behavior. Allowed values are flood and optimized_flooding. Default to flood.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dhcp_policy",
				Description: "(Optional) Map to provide dhcp_policy configurations.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dhcp_policy.name",
				Description: "(Optional) dhcp_policy name. Required if you specify the dhcp_policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DhcpPolicyName"),
			},
			{
				Name:        "dhcp_policy.version",
				Description: "(Optional) Version of dhcp_policy. Required if you specify the dhcp_policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DhcpPolicyVersion"),
			},
			{
				Name:        "dhcp_policy.dhcp_option_policy_name",
				Description: "(Optional) Name of dhcp_option_policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DhcpPolicyDhcpOptionPolicyName"),
			},
			{
				Name:        "dhcp_policy.dhcp_option_policy_version",
				Description: "(Optional) Version of dhcp_option_policy. Required if you specify the dhcp_policy.dhcp_option_policy_name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DhcpPolicyDhcpOptionPolicyVersion"),
			},
		},
	}
}

//// LIST FUNCTION
func listSchemaTemplateBd(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
			bdobjlist, err := curtemp.S("bds").Children()
			if err != nil {
				return nil, fmt.Errorf("Error getting Bd list: %v", err)
			}

			for _, curbd := range bdobjlist {
				bdobj := &SchemaTemplateBd{}
				bdobj.SchemaId = client.StripQuotes(curschema.S("id").String())
				bdobj.TemplateName = client.StripQuotes(curtemp.S("name").String())
				bdobj.Name = client.StripQuotes(curbd.S("name").String())
				bdobj.DisplayName = client.StripQuotes(curbd.S("displayName").String())
				if client.StripQuotes(curbd.S("vrfRef").String()) != "" {
					vrfinfo := strings.Split(client.StripQuotes(curbd.S("vrfRef").String()), "/")
					bdobj.VrfName = vrfinfo[6]
					bdobj.VrfSchemaId = vrfinfo[2]
					bdobj.VrfTemplateName = vrfinfo[4]
				}
				bdobj.Layer2UnknownUnicast = client.StripQuotes(curbd.S("l2UnknownUnicast").String())
				bdobj.IntersiteBumTraffic = client.StripQuotes(curbd.S("intersiteBumTrafficAllow").String())
				bdobj.OptimizeWanBandwidth = client.StripQuotes(curbd.S("optimizeWanBandwidth").String())
				bdobj.Layer2Stretch = client.StripQuotes(curbd.S("l2Stretch").String())
				bdobj.Layer3Multicast = client.StripQuotes(curbd.S("l3MCast").String())
				bdobj.ArpFlooding = client.StripQuotes(curbd.S("arpFlood").String())
				//bdobj.VirtualMacAddress = client.StripQuotes(curbd.S("virtualmac").String())
				bdobj.UnicastRouting = client.StripQuotes(curbd.S("unicastRouting").String())
				bdobj.Ipv6UnknownMulticastFlooding = client.StripQuotes(curbd.S("v6unkMcastAct").String())
				switch client.StripQuotes(curbd.S("multiDstPktAct").String()) {
				case "bd-flood":
					bdobj.MultiDestinationFlooding = "flood_in_bd"
				case "drop":
					bdobj.MultiDestinationFlooding = "drop"
				case "encap-flood":
					bdobj.MultiDestinationFlooding = "flood_in_encap"
				}

				bdobj.UnknownMulticastFlooding = client.StripQuotes(curbd.S("unkMcastAct").String())
				bdobj.Id = bdobj.SchemaId + "/template/" + bdobj.TemplateName + "/bd/" + bdobj.Name
				log.Printf("[TRACE] Record object: %s ", bdobj)
				d.StreamListItem(ctx, bdobj)
			}
		}
	}

	return nil, nil
}

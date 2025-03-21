package service

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/tools"
)

// OSPFInterfaceResourceModel describes the resource data model.
type QuaggaOSPFInterfaceResourceModel struct {
	Enabled       types.Bool   `tfsdk:"enabled"`
	InterfaceName types.String `tfsdk:"interfacename"`
	AuthType      types.String `tfsdk:"authtype"`
	AuthKey       types.Int64  `tfsdk:"authkey"`
	AuthKeyID     types.Int64  `tfsdk:"authkey_id"`
	Area          types.String `tfsdk:"area"`
	Cost          types.Int64  `tfsdk:"cost"`
	CostDemoted   types.Int64  `tfsdk:"cost_demoted"`
	// CarpDependOn       types.String `tfsdk:"carp_depend_on"`
	HelloInterval      types.Int64  `tfsdk:"hellointerval"`
	DeadInterval       types.Int64  `tfsdk:"deadinterval"`
	RetransmitInterval types.Int64  `tfsdk:"retransmitinterval"`
	RetransmitDelay    types.Int64  `tfsdk:"retransmitdelay"`
	TransmitDelay      types.Int64  `tfsdk:"transmitdelay"`
	Priority           types.Int64  `tfsdk:"priority"`
	BFD                types.Bool   `tfsdk:"bfd"`
	NetworkType        types.String `tfsdk:"networktype"`

	Id types.String `tfsdk:"id"`
}

func quaggaOSPFInterfaceResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure interface for OSPF",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this interface. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"interfacename": schema.StringAttribute{
				MarkdownDescription: "Select an interface where this settings apply. This uses an identifier like `lan` or or `opt2`. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"authtype": schema.StringAttribute{
				MarkdownDescription: "Defines security method for OSPF exchanges (None, plain, or MD5) to prevent unauthorized updates. Choose `MD5` or `plain`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf("", "message-digest", "plain"),
				},
			},
			"authkey": schema.Int64Attribute{
				MarkdownDescription: "Specifies a password or key used for plain or MD5 authentication. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"authkey_id": schema.Int64Attribute{
				MarkdownDescription: "Numeric identifier for MD5 authentication, ensuring correct key selection. The auth key ID. Defaults to `1`.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 255),
				},
			},
			"area": schema.StringAttribute{
				MarkdownDescription: "Assigns the network to an OSPF area using an identifier like 0.0.0.0 (Backbone Area). The Backbone Area connects other areas, supporting inter-area communication, while additional areas (e.g., 0.0.0.1, 0.0.0.255) segment the network logically to limit routing updates. Only use Area in Interface tab or in Network tab once. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"cost": schema.Int64Attribute{
				MarkdownDescription: "Sets the OSPF metric for path selection; lower costs are preferred paths within the area. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(40),
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"cost_demoted": schema.Int64Attribute{
				MarkdownDescription: "Specifies metric cost when interface is in backup mode via CARP, deprioritizing paths dynamically. Defaults to `65535`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(65535),
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			// "carp_depend_on": schema.StringAttribute{
			// 	MarkdownDescription: "Links the interface cost to a CARP VHID, adjusting costs based on primary or backup status.",
			// 	Optional:            true,
			// 	Computed:            true,
			// 	Default:             stringdefault.StaticString(""),
			// },
			"hellointerval": schema.Int64Attribute{
				MarkdownDescription: "Sets frequency (in seconds) of Hello packets to maintain OSPF neighbor relationships. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"deadinterval": schema.Int64Attribute{
				MarkdownDescription: "Defines the timeout period for OSPF neighbors; after this period, the neighbor is marked as down. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"retransmitinterval": schema.Int64Attribute{
				MarkdownDescription: "Time (seconds) to wait before resending Link-State Advertisements (LSAs) if acknowledgment is delayed. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"retransmitdelay": schema.Int64Attribute{
				MarkdownDescription: "Configures the hold time before LSAs are resent, accommodating slow or high-latency links. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"transmitdelay": schema.Int64Attribute{
				MarkdownDescription: "Configures the hold time before LSAs are resent, accommodating slow or high-latency links. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Determines the likelihood of becoming a Designated Router; higher values increase priority. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"bfd": schema.BoolAttribute{
				MarkdownDescription: "Activates Bidirectional Forwarding Detection for rapid link failure detection; peer configuration required. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"networktype": schema.StringAttribute{
				MarkdownDescription: "Defines the OSPF network type, impacting adjacency and LSA flooding methods. Defaults to `\"\"`.", // TODO: Add description for network types
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf("broadcast", "non-broadcast", "point-to-multipoint", "point-to-point"),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the interface.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func QuaggaOSPFInterfaceDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure interface for OSPF",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this interface. Defaults to `true`.",
				Computed:            true,
			},
			"interfacename": dschema.StringAttribute{
				MarkdownDescription: "Select an interface where this settings apply. This uses an identifier like `lan` or or `opt2`. Defaults to `\"\"`.",
				Computed:            true,
			},
			"authtype": dschema.StringAttribute{
				MarkdownDescription: "Defines security method for OSPF exchanges (None, plain, or MD5) to prevent unauthorized updates. Choose `MD5` or `plain`.",
				Computed:            true,
			},
			"authkey": dschema.Int64Attribute{
				MarkdownDescription: "Specifies a password or key used for plain or MD5 authentication. Defaults to `-1`.",
				Computed:            true,
			},
			"authkey_id": dschema.Int64Attribute{
				MarkdownDescription: "Numeric identifier for MD5 authentication, ensuring correct key selection. The auth key ID. Defaults to `1`.",
				Computed:            true,
			},
			"area": dschema.StringAttribute{
				MarkdownDescription: "Assigns the network to an OSPF area using an identifier like 0.0.0.0 (Backbone Area). The Backbone Area connects other areas, supporting inter-area communication, while additional areas (e.g., 0.0.0.1, 0.0.0.255) segment the network logically to limit routing updates. Only use Area in Interface tab or in Network tab once. Defaults to `\"\"`.",
				Computed:            true,
			},
			"cost": dschema.Int64Attribute{
				MarkdownDescription: "Sets the OSPF metric for path selection; lower costs are preferred paths within the area. Defaults to `-1`.",
				Computed:            true,
			},
			"cost_demoted": dschema.Int64Attribute{
				MarkdownDescription: "Specifies metric cost when interface is in backup mode via CARP, deprioritizing paths dynamically. Defaults to `65535`.",
				Computed:            true,
			},
			// "carp_depend_on": dschema.StringAttribute{
			// 	MarkdownDescription: "Links the interface cost to a CARP VHID, adjusting costs based on primary or backup status.",
			// 	Computed:            true,
			// },
			"hellointerval": dschema.Int64Attribute{
				MarkdownDescription: "Sets frequency (in seconds) of Hello packets to maintain OSPF neighbor relationships. Defaults to `-1`.",
				Computed:            true,
			},
			"deadinterval": dschema.Int64Attribute{
				MarkdownDescription: "Defines the timeout period for OSPF neighbors; after this period, the neighbor is marked as down. Defaults to `\"\"`.",
				Computed:            true,
			},
			"retransmitinterval": dschema.Int64Attribute{
				MarkdownDescription: "Time (seconds) to wait before resending Link-State Advertisements (LSAs) if acknowledgment is delayed. Defaults to `\"\"`.",
				Computed:            true,
			},
			"retransmitdelay": dschema.Int64Attribute{
				MarkdownDescription: "Configures the hold time before LSAs are resent, accommodating slow or high-latency links. Defaults to `\"\"`.",
				Computed:            true,
			},
			"transmitdelay": dschema.Int64Attribute{
				MarkdownDescription: "Configures the hold time before LSAs are resent, accommodating slow or high-latency links. Defaults to `\"\"`.",
				Computed:            true,
			},
			"priority": dschema.Int64Attribute{
				MarkdownDescription: "Determines the likelihood of becoming a Designated Router; higher values increase priority. Defaults to `\"\"`.",
				Computed:            true,
			},
			"bfd": dschema.BoolAttribute{
				MarkdownDescription: "Activates Bidirectional Forwarding Detection for rapid link failure detection; peer configuration required. Defaults to `false`.",
				Computed:            true,
			},
			"networktype": dschema.StringAttribute{
				MarkdownDescription: "Defines the OSPF network type, impacting adjacency and LSA flooding methods. Defaults to `\"\"`.", // TODO: Add description for network types
				Computed:            true,
			},
		},
	}
}

func convertQuaggaOSPFInterfaceSchemaToStruct(d *QuaggaOSPFInterfaceResourceModel) (*quagga.OSPFInterface, error) {
	return &quagga.OSPFInterface{
		Enabled:       tools.BoolToString(d.Enabled.ValueBool()),
		InterfaceName: api.SelectedMap(d.InterfaceName.ValueString()),
		AuthType:      api.SelectedMap(d.AuthType.ValueString()),
		AuthKey:       tools.Int64ToStringNegative(d.AuthKey.ValueInt64()),
		AuthKeyID:     tools.Int64ToString(d.AuthKeyID.ValueInt64()),
		Area:          d.Area.ValueString(),
		Cost:          tools.Int64ToString(d.Cost.ValueInt64()),
		CostDemoted:   tools.Int64ToString(d.CostDemoted.ValueInt64()),
		// CarpDependOn:       d.CarpDependOn.ValueString(),
		HelloInterval:      tools.Int64ToStringNegative(d.HelloInterval.ValueInt64()),
		DeadInterval:       tools.Int64ToStringNegative(d.DeadInterval.ValueInt64()),
		RetransmitInterval: tools.Int64ToStringNegative(d.RetransmitInterval.ValueInt64()),
		RetransmitDelay:    tools.Int64ToStringNegative(d.RetransmitDelay.ValueInt64()),
		TransmitDelay:      tools.Int64ToStringNegative(d.TransmitDelay.ValueInt64()),
		Priority:           tools.Int64ToStringNegative(d.Priority.ValueInt64()),
		BFD:                tools.BoolToString(d.BFD.ValueBool()),
		NetworkType:        api.SelectedMap(d.NetworkType.ValueString()),
	}, nil
}

func convertQuaggaOSPFInterfaceStructToSchema(d *quagga.OSPFInterface) (*QuaggaOSPFInterfaceResourceModel, error) {
	return &QuaggaOSPFInterfaceResourceModel{
		Enabled:       types.BoolValue(tools.StringToBool(d.Enabled)),
		InterfaceName: types.StringValue(d.InterfaceName.String()),
		AuthType:      types.StringValue(d.AuthType.String()),
		AuthKey:       types.Int64Value(tools.StringToInt64(d.AuthKey)),
		AuthKeyID:     types.Int64Value(tools.StringToInt64(d.AuthKeyID)),
		Area:          types.StringValue(d.Area),
		Cost:          types.Int64Value(tools.StringToInt64(d.Cost)),
		CostDemoted:   types.Int64Value(tools.StringToInt64(d.CostDemoted)),
		// CarpDependOn:       types.StringValue(d.CarpDependOn),
		HelloInterval:      types.Int64Value(tools.StringToInt64(d.HelloInterval)),
		DeadInterval:       types.Int64Value(tools.StringToInt64(d.DeadInterval)),
		RetransmitInterval: types.Int64Value(tools.StringToInt64(d.RetransmitInterval)),
		RetransmitDelay:    types.Int64Value(tools.StringToInt64(d.RetransmitDelay)),
		TransmitDelay:      types.Int64Value(tools.StringToInt64(d.TransmitDelay)),
		Priority:           types.Int64Value(tools.StringToInt64(d.Priority)),
		BFD:                types.BoolValue(tools.StringToBool(d.BFD)),
		NetworkType:        types.StringValue(d.NetworkType.String()),
	}, nil
}

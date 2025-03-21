package service

import (
	"context"
	"fmt"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &QuaggaOSPFInterfaceDataSource{}

func NewQuaggaOSPFInterfaceDataSource() datasource.DataSource {
	return &QuaggaOSPFInterfaceDataSource{}
}

// QuaggaOSPFInterfaceDataSource defines the data source implementation.
type QuaggaOSPFInterfaceDataSource struct {
	client opnsense.Client
}

func (d *QuaggaOSPFInterfaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quagga_ospf_interface"
}

func (d *QuaggaOSPFInterfaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = QuaggaOSPFInterfaceDataSourceSchema()
}

func (d *QuaggaOSPFInterfaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *opnsense.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = opnsense.NewClient(apiClient)
}

func (d *QuaggaOSPFInterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *QuaggaOSPFInterfaceResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get resource from OPNsense API
	resource, err := d.client.Quagga().GetOSPFInterface(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read neighbor, got error: %s", err))
		return
	}

	// Convert OPNsense struct to TF schema
	resourceModel, err := convertQuaggaOSPFInterfaceStructToSchema(resource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read neighbor, got error: %s", err))
		return
	}

	// ID cannot be added by convert... func, have to add here
	resourceModel.Id = data.Id

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}

package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/errs"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &QuaggaOSPFInterfaceResource{}
var _ resource.ResourceWithImportState = &QuaggaOSPFInterfaceResource{}

func NewQuaggaOSPFInterfaceResource() resource.Resource {
	return &QuaggaOSPFInterfaceResource{}
}

// QuaggaOSPFInterfaceResource defines the resource implementation.
type QuaggaOSPFInterfaceResource struct {
	client opnsense.Client
}

func (r *QuaggaOSPFInterfaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quagga_ospf_interface"
}

func (r *QuaggaOSPFInterfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = quaggaOSPFInterfaceResourceSchema()
}

func (r *QuaggaOSPFInterfaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = opnsense.NewClient(apiClient)
}

func (r *QuaggaOSPFInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *QuaggaOSPFInterfaceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	ospfInterface, err := convertQuaggaOSPFInterfaceSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ospf interface, got error: %s", err))
		return
	}

	// Add ospf interface to unbound
	id, err := r.client.Quagga().AddOSPFInterface(ctx, ospfInterface)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create ospf interface, got error: %s", err))
		return
	}

	// Tag new resource with ID from OPNsense
	data.Id = types.StringValue(id)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QuaggaOSPFInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *QuaggaOSPFInterfaceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ospf interface from OPNsense unbound API
	ospfInterface, err := r.client.Quagga().GetOSPFInterface(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, fmt.Sprintf("ospf interface not present in remote, removing from state"))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ospf interface, got error: %s", err))
		return
	}

	// Convert OPNsense struct to TF schema
	ospfInterfaceModel, err := convertQuaggaOSPFInterfaceStructToSchema(ospfInterface)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ospf interface, got error: %s", err))
		return
	}

	// ID cannot be added by convert... func, have to add here
	ospfInterfaceModel.Id = data.Id

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &ospfInterfaceModel)...)
}

func (r *QuaggaOSPFInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *QuaggaOSPFInterfaceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	ospfInterface, err := convertQuaggaOSPFInterfaceSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ospf interface, got error: %s", err))
		return
	}

	// Update ospf interface in unbound
	err = r.client.Quagga().UpdateOSPFInterface(ctx, data.Id.ValueString(), ospfInterface)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create ospf interface, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QuaggaOSPFInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *QuaggaOSPFInterfaceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Quagga().DeleteOSPFInterface(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete ospf interface, got error: %s", err))
		return
	}
}

func (r *QuaggaOSPFInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	turingpi "github.com/jlec/terraform-provider-turing-pi-bmc/internal/api"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &usbResource{}
	_ resource.ResourceWithImportState = &usbResource{}
)

func NewUsbResource() resource.Resource { //nolint:ireturn
	return &usbResource{}
}

// usbResource defines the resource implementation.
type usbResource usbDataSource

// type usbResource struct {
// 	client *turingpi.Client
// }

// usbResourceModel describes the data source data model.
type usbResourceModel usbDataSourceModel

// type usbResourceModel struct {
// 	ID   types.String `tfsdk:"id"`
// 	Mode types.Int64  `tfsdk:"mode"`
// 	Node types.Int64  `tfsdk:"node"`
// }

func (r *usbResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_usb"
}

func (r *usbResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Interface to the USB CM4 port in the Turing PI 2 board. It allows switching between the host " +
			"and device mode as well as mapping to one of the 4 nodes.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this resource",
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "USB port mode ('host' or 'device')",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("host", "device"),
				},
			},
			"node": schema.Int64Attribute{
				MarkdownDescription: "Node which USB port is mapped to",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, NodeCount),
				},
			},
		},
	}
}

func (r *usbResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	var ok bool

	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.client, ok = req.ProviderData.(*turingpi.Client)
	if !ok {
		resp.Diagnostics.AddError("Client Error", "failed to get client")
	}
}

func (r *usbResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan *usbResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	mode, err := turingpi.ModeToAPI(plan.Mode.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Api Error",
			fmt.Sprintf("Unable to convert API response, got error: %s", err),
		)

		return
	}

	_, err = r.client.SetUsb(ctx, mode, plan.Node.ValueInt64()-1)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read usb data, got error: %s", err),
		)

		return
	}

	plan.ID = types.StringValue("usb")
	plan.Mode = types.StringValue(plan.Mode.ValueString())
	plan.Node = types.Int64Value(plan.Node.ValueInt64())

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *usbResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan usbResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get USB state from BMC
	usb, err := r.client.GetUsb(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read usb data, got error: %s", err),
		)

		return
	}

	mode, err := turingpi.APIToMode(usb.Mode)
	if err != nil {
		resp.Diagnostics.AddError(
			"Api Error",
			fmt.Sprintf("Unable to convert API response, got error: %s", err),
		)

		return
	}

	plan.ID = types.StringValue("usb")
	plan.Mode = types.StringValue(mode)
	plan.Node = types.Int64Value(usb.Node + 1)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *usbResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan *usbResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	mode, err := turingpi.ModeToAPI(plan.Mode.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Api Error",
			fmt.Sprintf("Unable to convert API response, got error: %s", err),
		)

		return
	}

	_, err = r.client.SetUsb(ctx, mode, plan.Node.ValueInt64()-1)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to set usb data, got error: %s", err),
		)

		return
	}

	usb, err := r.client.GetUsb(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read usb data, got error: %s", err),
		)

		return
	}

	modeAPI, err := turingpi.APIToMode(usb.Mode)
	if err != nil {
		resp.Diagnostics.AddError(
			"Api Error",
			fmt.Sprintf("Unable to convert API response, got error: %s", err),
		)

		return
	}

	plan.ID = types.StringValue("usb")
	plan.Mode = types.StringValue(modeAPI)
	plan.Node = types.Int64Value(usb.Node + 1)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *usbResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data *usbResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *usbResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

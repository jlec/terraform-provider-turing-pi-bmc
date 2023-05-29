package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	turingpi "github.com/jlec/terraform-provider-turing-pi-bmc/internal/api"
)

// Ensure provider defined types fully satisfy framework interfaces.
// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &usbDataSource{}
	_ datasource.DataSourceWithConfigure = &usbDataSource{}
)

// NewUsbDataSource is a helper function to simplify the provider implementation.
func NewUsbDataSource() datasource.DataSource {
	return &usbDataSource{}
}

// usbDataSource defines the data source implementation.
type usbDataSource struct {
	client *turingpi.Client
}

// usbDataSourceModel describes the data source data model.
type usbDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Mode types.String `tfsdk:"mode"`
	Node types.Int64  `tfsdk:"node"`
}

func (d *usbDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_usb"
}

func (d *usbDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Turing PI Usb Data Source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID",
				Computed:    true,
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "USB mode",
				Computed:            true,
			},
			"node": schema.Int64Attribute{
				MarkdownDescription: "Node using USB",
				Computed:            true,
			},
		},
	}
}

// FIXME: RO attribute.
func (d *usbDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	var ok bool
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	d.client, ok = req.ProviderData.(*turingpi.Client)
	if !ok {
		resp.Diagnostics.AddError("Client Error", "failed to get client")
	}
}

func (d *usbDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data usbDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	usb, err := d.client.GetUsb()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read usb data, got error: %s", err))

		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	mode, err := turingpi.ApiToMode(usb.Mode)
	if err != nil {
		resp.Diagnostics.AddError("Api Error", fmt.Sprintf("Unable to convert API response, got error: %s", err))

		return
	}

	data.ID = types.StringValue("usb")
	data.Mode = types.StringValue(mode)
	data.Node = types.Int64Value(usb.Node + 1)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

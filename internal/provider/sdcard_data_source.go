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
	_ datasource.DataSource              = &sdCardDataSource{}
	_ datasource.DataSourceWithConfigure = &sdCardDataSource{}
)

func NewsdCardDataSource() datasource.DataSource {
	return &sdCardDataSource{}
}

// NewSDCardDataSource is a helper function to simplify the provider implementation.
func NewSDCardDataSource() datasource.DataSource {
	return &sdCardDataSource{}
}

// sdCardDataSource defines the data source implementation.
type sdCardDataSource struct {
	client *turingpi.Client
}

// sdCardDataSourceModel describes the data source data model.
type sdCardDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Total types.Int64  `tfsdk:"total"`
	Free  types.Int64  `tfsdk:"free"`
	Use   types.Int64  `tfsdk:"use"`
}

func (d *sdCardDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdcard"
}

func (d *sdCardDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Turing PI SDCard Data Source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID",
				Computed:    true,
			},
			"total": schema.Int64Attribute{
				MarkdownDescription: "Total capacity of SDCard",
				Computed:            true,
			},
			"free": schema.Int64Attribute{
				MarkdownDescription: "Total capacity of SDCard",
				Computed:            true,
			},
			"use": schema.Int64Attribute{
				MarkdownDescription: "Total capacity of SDCard",
				Computed:            true,
			},
		},
	}
}

// FIXME: RO attribute.
func (d *sdCardDataSource) Configure(
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

func (d *sdCardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data sdCardDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	sdCard, err := d.client.GetSDCard()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read sdcard data, got error: %s", err))

		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.ID = types.StringValue("sdcard")
	data.Total = types.Int64Value(sdCard.Total)
	data.Free = types.Int64Value(sdCard.Free)
	data.Use = types.Int64Value(sdCard.Use)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

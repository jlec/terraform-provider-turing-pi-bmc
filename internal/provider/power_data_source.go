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
	_ datasource.DataSource              = &powerDataSource{}
	_ datasource.DataSourceWithConfigure = &powerDataSource{}
)

// NewPowerDataSource is a helper function to simplify the provider implementation.
func NewPowerDataSource() datasource.DataSource {
	return &powerDataSource{}
}

// powerDataSource defines the data source implementation.
type powerDataSource struct {
	client *turingpi.Client
}

// powerDataSourceModel describes the data source data model.
type powerDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Node1 types.Int64  `tfsdk:"node1"`
	Node2 types.Int64  `tfsdk:"node2"`
	Node3 types.Int64  `tfsdk:"node3"`
	Node4 types.Int64  `tfsdk:"node4"`
}

func (d *powerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_power"
}

func (d *powerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Turing PI Power Data Source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID",
				Computed:    true,
			},
			"node1": schema.Int64Attribute{
				MarkdownDescription: "Power state of Node 1",
				Computed:            true,
			},
			"node2": schema.Int64Attribute{
				MarkdownDescription: "Power state of Node 2",
				Computed:            true,
			},
			"node3": schema.Int64Attribute{
				MarkdownDescription: "Power state of Node 3",
				Computed:            true,
			},
			"node4": schema.Int64Attribute{
				MarkdownDescription: "Power state of Node 4",
				Computed:            true,
			},
		},
	}
}

// FIXME: RO attribute.
func (d *powerDataSource) Configure(
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

func (d *powerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data powerDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	power, err := d.client.GetPower()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read power data, got error: %s", err))

		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.ID = types.StringValue("power")
	data.Node1 = types.Int64Value(power.Node1)
	data.Node2 = types.Int64Value(power.Node2)
	data.Node3 = types.Int64Value(power.Node3)
	data.Node4 = types.Int64Value(power.Node4)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

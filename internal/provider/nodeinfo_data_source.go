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
	_ datasource.DataSource              = &nodeInfoDataSource{}
	_ datasource.DataSourceWithConfigure = &nodeInfoDataSource{}
)

// NewNodeInfoDataSource is a helper function to simplify the provider implementation.
func NewNodeInfoDataSource() datasource.DataSource { //nolint:ireturn
	return &nodeInfoDataSource{}
}

// nodeInfoDataSource defines the data source implementation.
type nodeInfoDataSource struct {
	client *turingpi.Client
}

// nodeInfoDataSourceModel describes the data source data model.
type nodeInfoDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Node1 types.String `tfsdk:"node1"`
	Node2 types.String `tfsdk:"node2"`
	Node3 types.String `tfsdk:"node3"`
	Node4 types.String `tfsdk:"node4"`
}

func (d *nodeInfoDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_nodeinfo"
}

func (d *nodeInfoDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Turing PI NodeInfo Data Source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID",
				Computed:    true,
			},
			"node1": schema.StringAttribute{
				MarkdownDescription: "Information for Node 1",
				Computed:            true,
			},
			"node2": schema.StringAttribute{
				MarkdownDescription: "Information for Node 2",
				Computed:            true,
			},
			"node3": schema.StringAttribute{
				MarkdownDescription: "Information for Node 3",
				Computed:            true,
			},
			"node4": schema.StringAttribute{
				MarkdownDescription: "Information for Node 4",
				Computed:            true,
			},
		},
	}
}

// FIXME: RO attribute.
func (d *nodeInfoDataSource) Configure(
	_ context.Context,
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

func (d *nodeInfoDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data nodeInfoDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	nodeInfo, err := d.client.GetNodeInfo(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read nodeinfo data, got error: %s", err),
		)

		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.ID = types.StringValue("nodeinfo")
	data.Node1 = types.StringValue(nodeInfo.Node1)
	data.Node2 = types.StringValue(nodeInfo.Node2)
	data.Node3 = types.StringValue(nodeInfo.Node3)
	data.Node4 = types.StringValue(nodeInfo.Node4)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

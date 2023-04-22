package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	turingpi "github.com/jlec/terraform-provider-turing-pi-bmc/internal/api"
)

// Ensure TuringPiBMCProvider satisfies various provider interfaces.
var _ provider.Provider = &TuringPiBMCProvider{}

// TuringPiBMCProvider defines the provider implementation.
type TuringPiBMCProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// TuringPiBMCProviderModel describes the provider data model.
type TuringPiBMCProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *TuringPiBMCProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "turing-pi-bmc"
	resp.Version = p.version
}

func (p *TuringPiBMCProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Turing Pi BMC endpoint",
				Required:            true,
			},
		},
	}
}

func (p *TuringPiBMCProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data TuringPiBMCProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }
	if data.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Turing Pi API Host",
			"You need to provide the IP or DNS name pointing to your Turing PI BMC",
		)
	}
	// Example client configuration for data sources and resources
	client, err := turingpi.NewClient(data.Endpoint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Something went wrong, got error: %s", err))

		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *TuringPiBMCProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// NewSDCardDataSource,
	}
}

func (p *TuringPiBMCProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewNodeInfoDataSource,
		NewPowerDataSource,
		NewSDCardDataSource,
		NewUsbDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TuringPiBMCProvider{
			version: version,
		}
	}
}

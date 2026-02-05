package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the provider satisfies the provider.Provider interface.
var _ provider.Provider = &CidrguardProvider{}

// CidrguardProvider is the provider implementation.
type CidrguardProvider struct {
	version string
}

// New returns a new provider.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CidrguardProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *CidrguardProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cidrguard"
	resp.Version = p.version
}

// Schema defines the provider-level schema.
func (p *CidrguardProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Configure prepares a provider for data source and resource requests.
func (p *CidrguardProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

// DataSources defines the data sources implemented in the provider.
func (p *CidrguardProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCidrguardRegistryDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *CidrguardProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

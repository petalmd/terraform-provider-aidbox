// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"net/http"
	"os" // Import for environment variables
	"terraform-provider-aidbox/internal/aidbox"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure AidboxProvider satisfies various provider interfaces.
var _ provider.Provider = &AidboxProvider{}
var _ provider.ProviderWithFunctions = &AidboxProvider{}

type AidboxProvider struct {
	version string
}

type AidboxProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

type Client interface {
	CreateLicense(cxt context.Context, name, product, licenseType string) (aidbox.LicenseResponse, error)
	GetLicense(ctx context.Context, licenseID string) (aidbox.LicenseResponse, error)
	DeleteLicense(ctx context.Context, licenseID string) error
}

// This structure holds the configuration data which can be used across resources
type ProviderData struct {
	Endpoint string
	Token    string
	Client   Client
}

func (p *AidboxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "aidbox"
	resp.Version = p.version
}

func (p *AidboxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Aidbox RPC API endpoint",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Aidbox API token",
				Optional:            true,
			},
		},
	}
}

func (p *AidboxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AidboxProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Set default endpoint if not provided
	if data.Endpoint.IsNull() || data.Endpoint.IsUnknown() || data.Endpoint.ValueString() == "" {
		defaultEndpoint := basetypes.NewStringValue("https://aidbox.app/rpc")
		data.Endpoint = defaultEndpoint
	}

	// Handle token; get from environment variable if not provided
	if data.Token.IsNull() || data.Token.IsUnknown() || data.Token.ValueString() == "" {
		tokenEnv := os.Getenv("AIDBOX_API_TOKEN")
		if tokenEnv != "" {
			data.Token = basetypes.NewStringValue(tokenEnv)
		} else {
			resp.Diagnostics.AddError(
				"No API Token Provided",
				"Please provide a 'token' in the provider configuration or through the 'AIDBOX_API_TOKEN' environment variable.",
			)
			return
		}
	}

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = &ProviderData{
		Endpoint: data.Endpoint.ValueString(),
		Token:    data.Token.ValueString(),
		Client:   aidbox.NewClient(data.Endpoint.ValueString(), data.Token.ValueString()),
	}
}

func (p *AidboxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewLicenseResource,
	}
}

func (p *AidboxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *AidboxProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AidboxProvider{
			version: version,
		}
	}
}

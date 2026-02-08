// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/ackack-io/terraform-provider-ackack/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure AckackProvider satisfies various provider interfaces.
var _ provider.Provider = &AckackProvider{}

// AckackProvider defines the provider implementation.
type AckackProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// AckackProviderModel describes the provider data model.
type AckackProviderModel struct {
	APIKey   types.String `tfsdk:"api_key"`
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *AckackProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ackack"
	resp.Version = p.version
}

func (p *AckackProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The ackack provider allows you to manage uptime monitors, alerts, systems, and reports on ackack.io.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key for authenticating with ackack.io. Can also be set via the `ACKACK_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The ackack.io API endpoint. Defaults to `https://api.ackack.io`. Can also be set via the `ACKACK_ENDPOINT` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *AckackProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AckackProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check environment variables first, then use config values
	apiKey := os.Getenv("ACKACK_API_KEY")
	if !data.APIKey.IsNull() {
		apiKey = data.APIKey.ValueString()
	}

	endpoint := os.Getenv("ACKACK_ENDPOINT")
	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}

	// Validate required configuration
	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing ackack API Key",
			"The provider cannot create the ackack API client as there is a missing or empty value for the ackack API key. "+
				"Set the api_key value in the configuration or use the ACKACK_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		return
	}

	// Create API client
	c, err := client.NewClient(apiKey, endpoint, p.version)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create ackack API Client",
			"An unexpected error occurred when creating the ackack API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *AckackProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewMonitorResource,
		NewAlertResource,
		NewSystemResource,
		NewReportResource,
	}
}

func (p *AckackProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMonitorDataSource,
		NewMonitorsDataSource,
		NewAlertDataSource,
		NewAlertsDataSource,
		NewSystemDataSource,
		NewSystemsDataSource,
		NewMonitorResultsDataSource,
		NewMonitorUptimeDataSource,
		NewMonitorIncidentsDataSource,
		NewMonitorHealthDataSource,
		NewNotificationsDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AckackProvider{
			version: version,
		}
	}
}

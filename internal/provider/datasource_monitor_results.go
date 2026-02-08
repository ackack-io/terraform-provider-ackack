// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/ackack-io/terraform-provider-ackack/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MonitorResultsDataSource{}

func NewMonitorResultsDataSource() datasource.DataSource {
	return &MonitorResultsDataSource{}
}

// MonitorResultsDataSource defines the data source implementation.
type MonitorResultsDataSource struct {
	client *client.Client
}

// MonitorResultsDataSourceModel describes the data source data model.
type MonitorResultsDataSourceModel struct {
	MonitorID types.String               `tfsdk:"monitor_id"`
	Limit     types.Int64                `tfsdk:"limit"`
	Results   []MonitorResultItemModel   `tfsdk:"results"`
}

// MonitorResultItemModel describes a single check result.
type MonitorResultItemModel struct {
	ID                        types.Int64  `tfsdk:"id"`
	Status                    types.String `tfsdk:"status"`
	ResponseTime              types.Int64  `tfsdk:"response_time"`
	ResponseSizeBytes         types.Int64  `tfsdk:"response_size_bytes"`
	Timestamp                 types.String `tfsdk:"timestamp"`
	Region                    types.String `tfsdk:"region"`
	Message                   types.String `tfsdk:"message"`
	ErrorType                 types.String `tfsdk:"error_type"`
	StatusCode                types.Int64  `tfsdk:"status_code"`
	DNSResponse               types.String `tfsdk:"dns_response"`
	TLSVersion                types.String `tfsdk:"tls_version"`
	CertificateExpirationDays types.Int64  `tfsdk:"certificate_expiration_days"`
}

func (d *MonitorResultsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor_results"
}

func (d *MonitorResultsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get recent check results for a monitor.",

		Attributes: map[string]schema.Attribute{
			"monitor_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the monitor.",
				Required:            true,
			},
			"limit": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of results to return. Default is 100, max is 1000.",
				Optional:            true,
			},
			"results": schema.ListNestedAttribute{
				MarkdownDescription: "List of check results.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "The result ID.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "The check status.",
							Computed:            true,
						},
						"response_time": schema.Int64Attribute{
							MarkdownDescription: "Response time in milliseconds.",
							Computed:            true,
						},
						"response_size_bytes": schema.Int64Attribute{
							MarkdownDescription: "Response size in bytes.",
							Computed:            true,
						},
						"timestamp": schema.StringAttribute{
							MarkdownDescription: "The timestamp of the check.",
							Computed:            true,
						},
						"region": schema.StringAttribute{
							MarkdownDescription: "The region where the check was performed.",
							Computed:            true,
						},
						"message": schema.StringAttribute{
							MarkdownDescription: "Any message associated with the check.",
							Computed:            true,
						},
						"error_type": schema.StringAttribute{
							MarkdownDescription: "The type of error if the check failed.",
							Computed:            true,
						},
						"status_code": schema.Int64Attribute{
							MarkdownDescription: "HTTP status code (for HTTP monitors).",
							Computed:            true,
						},
						"dns_response": schema.StringAttribute{
							MarkdownDescription: "DNS response (for DNS monitors).",
							Computed:            true,
						},
						"tls_version": schema.StringAttribute{
							MarkdownDescription: "TLS version (for SSL monitors).",
							Computed:            true,
						},
						"certificate_expiration_days": schema.Int64Attribute{
							MarkdownDescription: "Days until certificate expiration (for SSL monitors).",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *MonitorResultsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *MonitorResultsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MonitorResultsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	limit := 0
	if !data.Limit.IsNull() {
		limit = int(data.Limit.ValueInt64())
	}

	results, err := d.client.GetMonitorResults(ctx, data.MonitorID.ValueString(), limit)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get monitor results, got error: %s", err))
		return
	}

	data.Results = make([]MonitorResultItemModel, len(results))
	for i, result := range results {
		data.Results[i] = MonitorResultItemModel{
			ID:                types.Int64Value(int64(result.ID)),
			Status:            types.StringValue(result.Status),
			ResponseTime:      types.Int64Value(int64(result.ResponseTime)),
			ResponseSizeBytes: types.Int64Value(int64(result.ResponseSizeBytes)),
			Timestamp:         types.StringValue(result.Timestamp),
		}
		if result.Region != "" {
			data.Results[i].Region = types.StringValue(result.Region)
		}
		if result.Message != "" {
			data.Results[i].Message = types.StringValue(result.Message)
		}
		if result.ErrorType != "" {
			data.Results[i].ErrorType = types.StringValue(result.ErrorType)
		}
		if result.StatusCode != 0 {
			data.Results[i].StatusCode = types.Int64Value(int64(result.StatusCode))
		}
		if result.DNSResponse != "" {
			data.Results[i].DNSResponse = types.StringValue(result.DNSResponse)
		}
		if result.TLSVersion != "" {
			data.Results[i].TLSVersion = types.StringValue(result.TLSVersion)
		}
		if result.CertificateExpirationDays != 0 {
			data.Results[i].CertificateExpirationDays = types.Int64Value(int64(result.CertificateExpirationDays))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

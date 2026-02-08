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
var _ datasource.DataSource = &MonitorDataSource{}

func NewMonitorDataSource() datasource.DataSource {
	return &MonitorDataSource{}
}

// MonitorDataSource defines the data source implementation.
type MonitorDataSource struct {
	client *client.Client
}

// MonitorDataSourceModel describes the data source data model.
type MonitorDataSourceModel struct {
	ID               types.String  `tfsdk:"id"`
	Name             types.String  `tfsdk:"name"`
	Type             types.String  `tfsdk:"type"`
	IsEnabled        types.Bool    `tfsdk:"is_enabled"`
	FrequencySeconds types.Int64   `tfsdk:"frequency_seconds"`
	TimeoutMs        types.Int64   `tfsdk:"timeout_ms"`
	Retries          types.Int64   `tfsdk:"retries"`
	GeneralRegion    types.String  `tfsdk:"general_region"`
	SpecificRegion   types.String  `tfsdk:"specific_region"`
	Status           types.String  `tfsdk:"status"`
	UptimePercentage types.Float64 `tfsdk:"uptime_percentage"`
	LastChecked      types.String  `tfsdk:"last_checked"`
	CreatedAt        types.String  `tfsdk:"created_at"`
	UpdatedAt        types.String  `tfsdk:"updated_at"`

	// HTTP specific
	URL                types.String `tfsdk:"url"`
	ExpectedStatusCode types.Int64  `tfsdk:"expected_status_code"`
	ValidateStatus     types.Bool   `tfsdk:"validate_status"`
	ValidateBody       types.Bool   `tfsdk:"validate_body"`
	BodyPattern        types.String `tfsdk:"body_pattern"`
	Headers            types.String `tfsdk:"headers"`

	// DNS specific
	DNSRecordType types.String `tfsdk:"dns_record_type"`
	ExpectedValue types.String `tfsdk:"expected_value"`
	Nameserver    types.String `tfsdk:"nameserver"`

	// TCP specific
	Host types.String `tfsdk:"host"`
	Port types.Int64  `tfsdk:"port"`

	// SSL specific
	Domain                   types.String `tfsdk:"domain"`
	CheckExpirationThreshold types.Bool   `tfsdk:"check_expiration_threshold"`
	ExpirationThreshold      types.Int64  `tfsdk:"expiration_threshold"`
	CheckProtocolVersion     types.Bool   `tfsdk:"check_protocol_version"`
	MinimumProtocol          types.String `tfsdk:"minimum_protocol"`
}

func (d *MonitorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (d *MonitorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a specific monitor.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the monitor.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the monitor.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of monitor (http, dns, ssl, tcp).",
				Computed:            true,
			},
			"is_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the monitor is enabled.",
				Computed:            true,
			},
			"frequency_seconds": schema.Int64Attribute{
				MarkdownDescription: "How often the monitor checks, in seconds.",
				Computed:            true,
			},
			"timeout_ms": schema.Int64Attribute{
				MarkdownDescription: "Timeout for each check, in milliseconds.",
				Computed:            true,
			},
			"retries": schema.Int64Attribute{
				MarkdownDescription: "Number of retries before marking as failed.",
				Computed:            true,
			},
			"general_region": schema.StringAttribute{
				MarkdownDescription: "The general region for monitoring.",
				Computed:            true,
			},
			"specific_region": schema.StringAttribute{
				MarkdownDescription: "The specific region for monitoring.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The current status of the monitor.",
				Computed:            true,
			},
			"uptime_percentage": schema.Float64Attribute{
				MarkdownDescription: "The uptime percentage of the monitor.",
				Computed:            true,
			},
			"last_checked": schema.StringAttribute{
				MarkdownDescription: "The timestamp of the last check.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the monitor was created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the monitor was last updated.",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL to monitor (HTTP monitors).",
				Computed:            true,
			},
			"expected_status_code": schema.Int64Attribute{
				MarkdownDescription: "The expected HTTP status code.",
				Computed:            true,
			},
			"validate_status": schema.BoolAttribute{
				MarkdownDescription: "Whether to validate the HTTP status code.",
				Computed:            true,
			},
			"validate_body": schema.BoolAttribute{
				MarkdownDescription: "Whether to validate the response body.",
				Computed:            true,
			},
			"body_pattern": schema.StringAttribute{
				MarkdownDescription: "The pattern to match in the response body.",
				Computed:            true,
			},
			"headers": schema.StringAttribute{
				MarkdownDescription: "HTTP headers as a JSON string.",
				Computed:            true,
			},
			"dns_record_type": schema.StringAttribute{
				MarkdownDescription: "The DNS record type to query.",
				Computed:            true,
			},
			"expected_value": schema.StringAttribute{
				MarkdownDescription: "The expected DNS record value.",
				Computed:            true,
			},
			"nameserver": schema.StringAttribute{
				MarkdownDescription: "The nameserver to query.",
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The host to connect to (TCP monitors).",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port to connect to (TCP monitors).",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The domain to check SSL certificate for.",
				Computed:            true,
			},
			"check_expiration_threshold": schema.BoolAttribute{
				MarkdownDescription: "Whether to check if the certificate is expiring soon.",
				Computed:            true,
			},
			"expiration_threshold": schema.Int64Attribute{
				MarkdownDescription: "Days before expiration to alert.",
				Computed:            true,
			},
			"check_protocol_version": schema.BoolAttribute{
				MarkdownDescription: "Whether to check the TLS protocol version.",
				Computed:            true,
			},
			"minimum_protocol": schema.StringAttribute{
				MarkdownDescription: "The minimum TLS protocol version.",
				Computed:            true,
			},
		},
	}
}

func (d *MonitorDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MonitorDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := d.client.GetMonitor(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor, got error: %s", err))
		return
	}

	data.Name = types.StringValue(monitor.Name)
	data.Type = types.StringValue(monitor.Type)
	data.IsEnabled = types.BoolValue(monitor.IsEnabled)
	data.FrequencySeconds = types.Int64Value(int64(monitor.FrequencySeconds))
	data.TimeoutMs = types.Int64Value(int64(monitor.TimeoutMs))
	data.Retries = types.Int64Value(int64(monitor.Retries))
	data.Status = types.StringValue(monitor.Status)
	data.UptimePercentage = types.Float64Value(monitor.UptimePercentage)
	data.CreatedAt = types.StringValue(monitor.CreatedAt)
	data.UpdatedAt = types.StringValue(monitor.UpdatedAt)

	if monitor.GeneralRegion != "" {
		data.GeneralRegion = types.StringValue(monitor.GeneralRegion)
	}
	if monitor.SpecificRegion != "" {
		data.SpecificRegion = types.StringValue(monitor.SpecificRegion)
	}
	if monitor.LastChecked != "" {
		data.LastChecked = types.StringValue(monitor.LastChecked)
	}
	if monitor.URL != "" {
		data.URL = types.StringValue(monitor.URL)
	}
	if monitor.ExpectedStatusCode != 0 {
		data.ExpectedStatusCode = types.Int64Value(int64(monitor.ExpectedStatusCode))
	}
	data.ValidateStatus = types.BoolValue(monitor.ValidateStatus)
	data.ValidateBody = types.BoolValue(monitor.ValidateBody)
	if monitor.BodyPattern != "" {
		data.BodyPattern = types.StringValue(monitor.BodyPattern)
	}
	if monitor.Headers != "" {
		data.Headers = types.StringValue(monitor.Headers)
	}
	if monitor.DNSRecordType != "" {
		data.DNSRecordType = types.StringValue(monitor.DNSRecordType)
	}
	if monitor.ExpectedValue != "" {
		data.ExpectedValue = types.StringValue(monitor.ExpectedValue)
	}
	if monitor.Nameserver != "" {
		data.Nameserver = types.StringValue(monitor.Nameserver)
	}
	if monitor.Host != "" {
		data.Host = types.StringValue(monitor.Host)
	}
	if monitor.Port != 0 {
		data.Port = types.Int64Value(int64(monitor.Port))
	}
	if monitor.Domain != "" {
		data.Domain = types.StringValue(monitor.Domain)
	}
	data.CheckExpirationThreshold = types.BoolValue(monitor.CheckExpirationThreshold)
	if monitor.ExpirationThreshold != 0 {
		data.ExpirationThreshold = types.Int64Value(int64(monitor.ExpirationThreshold))
	}
	data.CheckProtocolVersion = types.BoolValue(monitor.CheckProtocolVersion)
	if monitor.MinimumProtocol != "" {
		data.MinimumProtocol = types.StringValue(monitor.MinimumProtocol)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

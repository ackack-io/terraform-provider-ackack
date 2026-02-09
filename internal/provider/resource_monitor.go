// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/ackack-io/terraform-provider-ackack/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MonitorResource{}
var _ resource.ResourceWithImportState = &MonitorResource{}

func NewMonitorResource() resource.Resource {
	return &MonitorResource{}
}

// MonitorResource defines the resource implementation.
type MonitorResource struct {
	client *client.Client
}

// MonitorResourceModel describes the resource data model.
type MonitorResourceModel struct {
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

func (r *MonitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (r *MonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an uptime monitor on ackack.io.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the monitor.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the monitor.",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of monitor. Must be one of: `http`, `dns`, `ssl`, `tcp`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("http", "dns", "ssl", "tcp"),
				},
			},
			"is_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the monitor is enabled. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"frequency_seconds": schema.Int64Attribute{
				MarkdownDescription: "How often to check the monitor, in seconds. Defaults to `60`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(60),
			},
			"timeout_ms": schema.Int64Attribute{
				MarkdownDescription: "Timeout for each check, in milliseconds. Defaults to `10000`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(10000),
			},
			"retries": schema.Int64Attribute{
				MarkdownDescription: "Number of retries before marking as failed.",
				Optional:            true,
				Computed:            true,
			},
			"general_region": schema.StringAttribute{
				MarkdownDescription: "The general region for monitoring (e.g., `us`, `eu`, `asia`).",
				Optional:            true,
				Computed:            true,
			},
			"specific_region": schema.StringAttribute{
				MarkdownDescription: "The specific region for monitoring.",
				Optional:            true,
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

			// HTTP specific
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL to monitor. Required for HTTP monitors.",
				Optional:            true,
			},
			"expected_status_code": schema.Int64Attribute{
				MarkdownDescription: "The expected HTTP status code. Defaults to `200`.",
				Optional:            true,
			},
			"validate_status": schema.BoolAttribute{
				MarkdownDescription: "Whether to validate the HTTP status code.",
				Optional:            true,
				Computed:            true,
			},
			"validate_body": schema.BoolAttribute{
				MarkdownDescription: "Whether to validate the response body.",
				Optional:            true,
				Computed:            true,
			},
			"body_pattern": schema.StringAttribute{
				MarkdownDescription: "The pattern to match in the response body.",
				Optional:            true,
			},
			"headers": schema.StringAttribute{
				MarkdownDescription: "HTTP headers as a JSON string.",
				Optional:            true,
			},

			// DNS specific
			"dns_record_type": schema.StringAttribute{
				MarkdownDescription: "The DNS record type to query (e.g., `A`, `AAAA`, `CNAME`). Required for DNS monitors.",
				Optional:            true,
			},
			"expected_value": schema.StringAttribute{
				MarkdownDescription: "The expected DNS record value.",
				Optional:            true,
			},
			"nameserver": schema.StringAttribute{
				MarkdownDescription: "The nameserver to query.",
				Optional:            true,
			},

			// TCP specific
			"host": schema.StringAttribute{
				MarkdownDescription: "The host to connect to. Required for TCP monitors.",
				Optional:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port to connect to. Required for TCP monitors.",
				Optional:            true,
			},

			// SSL specific
			"domain": schema.StringAttribute{
				MarkdownDescription: "The domain to check SSL certificate for. Required for SSL monitors.",
				Optional:            true,
			},
			"check_expiration_threshold": schema.BoolAttribute{
				MarkdownDescription: "Whether to check if the certificate is expiring soon.",
				Optional:            true,
				Computed:            true,
			},
			"expiration_threshold": schema.Int64Attribute{
				MarkdownDescription: "Days before expiration to alert.",
				Optional:            true,
			},
			"check_protocol_version": schema.BoolAttribute{
				MarkdownDescription: "Whether to check the TLS protocol version.",
				Optional:            true,
				Computed:            true,
			},
			"minimum_protocol": schema.StringAttribute{
				MarkdownDescription: "The minimum TLS protocol version (e.g., `TLS1.2`, `TLS1.3`).",
				Optional:            true,
			},
		},
	}
}

func (r *MonitorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *MonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MonitorResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := r.buildCreateRequest(&data)

	monitor, err := r.client.CreateMonitor(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create monitor, got error: %s", err))
		return
	}

	r.updateModelFromResponse(&data, monitor)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MonitorResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := r.client.GetMonitor(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor, got error: %s", err))
		return
	}

	r.updateModelFromResponse(&data, monitor)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MonitorResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := r.buildUpdateRequest(&data)

	monitor, err := r.client.UpdateMonitor(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update monitor, got error: %s", err))
		return
	}

	r.updateModelFromResponse(&data, monitor)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MonitorResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteMonitor(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete monitor, got error: %s", err))
		return
	}
}

func (r *MonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *MonitorResource) buildCreateRequest(data *MonitorResourceModel) client.CreateMonitorRequest {
	req := client.CreateMonitorRequest{
		Name: data.Name.ValueString(),
		Type: data.Type.ValueString(),
	}

	if !data.IsEnabled.IsNull() {
		isEnabled := data.IsEnabled.ValueBool()
		req.IsEnabled = &isEnabled
	}
	if !data.FrequencySeconds.IsNull() {
		req.FrequencySeconds = int(data.FrequencySeconds.ValueInt64())
	}
	if !data.TimeoutMs.IsNull() {
		req.TimeoutMs = int(data.TimeoutMs.ValueInt64())
	}
	if !data.Retries.IsNull() {
		req.Retries = int(data.Retries.ValueInt64())
	}
	if !data.GeneralRegion.IsNull() {
		req.GeneralRegion = data.GeneralRegion.ValueString()
	}
	if !data.SpecificRegion.IsNull() {
		req.SpecificRegion = data.SpecificRegion.ValueString()
	}

	// HTTP specific
	if !data.URL.IsNull() {
		req.URL = data.URL.ValueString()
	}
	if !data.ExpectedStatusCode.IsNull() {
		req.ExpectedStatusCode = int(data.ExpectedStatusCode.ValueInt64())
	}
	if !data.ValidateStatus.IsNull() {
		validateStatus := data.ValidateStatus.ValueBool()
		req.ValidateStatus = &validateStatus
	}
	if !data.ValidateBody.IsNull() {
		validateBody := data.ValidateBody.ValueBool()
		req.ValidateBody = &validateBody
	}
	if !data.BodyPattern.IsNull() {
		req.BodyPattern = data.BodyPattern.ValueString()
	}
	if !data.Headers.IsNull() {
		req.Headers = data.Headers.ValueString()
	}

	// DNS specific
	if !data.DNSRecordType.IsNull() {
		req.DNSRecordType = data.DNSRecordType.ValueString()
	}
	if !data.ExpectedValue.IsNull() {
		req.ExpectedValue = data.ExpectedValue.ValueString()
	}
	if !data.Nameserver.IsNull() {
		req.Nameserver = data.Nameserver.ValueString()
	}

	// TCP specific
	if !data.Host.IsNull() {
		req.Host = data.Host.ValueString()
	}
	if !data.Port.IsNull() {
		req.Port = int(data.Port.ValueInt64())
	}

	// SSL specific
	if !data.Domain.IsNull() {
		req.Domain = data.Domain.ValueString()
	}
	if !data.CheckExpirationThreshold.IsNull() {
		checkExp := data.CheckExpirationThreshold.ValueBool()
		req.CheckExpirationThreshold = &checkExp
	}
	if !data.ExpirationThreshold.IsNull() {
		req.ExpirationThreshold = int(data.ExpirationThreshold.ValueInt64())
	}
	if !data.CheckProtocolVersion.IsNull() {
		checkProto := data.CheckProtocolVersion.ValueBool()
		req.CheckProtocolVersion = &checkProto
	}
	if !data.MinimumProtocol.IsNull() {
		req.MinimumProtocol = data.MinimumProtocol.ValueString()
	}

	return req
}

func (r *MonitorResource) buildUpdateRequest(data *MonitorResourceModel) client.UpdateMonitorRequest {
	req := client.UpdateMonitorRequest{
		Name: data.Name.ValueString(),
		Type: data.Type.ValueString(),
	}

	if !data.IsEnabled.IsNull() {
		isEnabled := data.IsEnabled.ValueBool()
		req.IsEnabled = &isEnabled
	}
	if !data.FrequencySeconds.IsNull() {
		req.FrequencySeconds = int(data.FrequencySeconds.ValueInt64())
	}
	if !data.TimeoutMs.IsNull() {
		req.TimeoutMs = int(data.TimeoutMs.ValueInt64())
	}
	if !data.Retries.IsNull() {
		req.Retries = int(data.Retries.ValueInt64())
	}
	if !data.GeneralRegion.IsNull() {
		req.GeneralRegion = data.GeneralRegion.ValueString()
	}
	if !data.SpecificRegion.IsNull() {
		req.SpecificRegion = data.SpecificRegion.ValueString()
	}

	// HTTP specific
	if !data.URL.IsNull() {
		req.URL = data.URL.ValueString()
	}
	if !data.ExpectedStatusCode.IsNull() {
		req.ExpectedStatusCode = int(data.ExpectedStatusCode.ValueInt64())
	}
	if !data.ValidateStatus.IsNull() {
		validateStatus := data.ValidateStatus.ValueBool()
		req.ValidateStatus = &validateStatus
	}
	if !data.ValidateBody.IsNull() {
		validateBody := data.ValidateBody.ValueBool()
		req.ValidateBody = &validateBody
	}
	if !data.BodyPattern.IsNull() {
		req.BodyPattern = data.BodyPattern.ValueString()
	}
	if !data.Headers.IsNull() {
		req.Headers = data.Headers.ValueString()
	}

	// DNS specific
	if !data.DNSRecordType.IsNull() {
		req.DNSRecordType = data.DNSRecordType.ValueString()
	}
	if !data.ExpectedValue.IsNull() {
		req.ExpectedValue = data.ExpectedValue.ValueString()
	}
	if !data.Nameserver.IsNull() {
		req.Nameserver = data.Nameserver.ValueString()
	}

	// TCP specific
	if !data.Host.IsNull() {
		req.Host = data.Host.ValueString()
	}
	if !data.Port.IsNull() {
		req.Port = int(data.Port.ValueInt64())
	}

	// SSL specific
	if !data.Domain.IsNull() {
		req.Domain = data.Domain.ValueString()
	}
	if !data.CheckExpirationThreshold.IsNull() {
		checkExp := data.CheckExpirationThreshold.ValueBool()
		req.CheckExpirationThreshold = &checkExp
	}
	if !data.ExpirationThreshold.IsNull() {
		req.ExpirationThreshold = int(data.ExpirationThreshold.ValueInt64())
	}
	if !data.CheckProtocolVersion.IsNull() {
		checkProto := data.CheckProtocolVersion.ValueBool()
		req.CheckProtocolVersion = &checkProto
	}
	if !data.MinimumProtocol.IsNull() {
		req.MinimumProtocol = data.MinimumProtocol.ValueString()
	}

	return req
}

// normalizeTimestamp parses a timestamp and re-formats it with microsecond
// precision so that values stored in state always match what the API returns.
func normalizeTimestamp(ts string) string {
	t, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		return ts
	}
	return t.Truncate(time.Microsecond).Format("2006-01-02T15:04:05.999999Z07:00")
}

func (r *MonitorResource) updateModelFromResponse(data *MonitorResourceModel, monitor *client.Monitor) {
	data.ID = types.StringValue(monitor.ID)
	data.Name = types.StringValue(monitor.Name)
	data.Type = types.StringValue(monitor.Type)
	data.IsEnabled = types.BoolValue(monitor.IsEnabled)
	data.FrequencySeconds = types.Int64Value(int64(monitor.FrequencySeconds))
	data.TimeoutMs = types.Int64Value(int64(monitor.TimeoutMs))
	data.Retries = types.Int64Value(int64(monitor.Retries))
	data.Status = types.StringValue(monitor.Status)
	data.UptimePercentage = types.Float64Value(monitor.UptimePercentage)
	data.CreatedAt = types.StringValue(normalizeTimestamp(monitor.CreatedAt))
	data.UpdatedAt = types.StringValue(normalizeTimestamp(monitor.UpdatedAt))

	// Set optional string fields - use null if empty to ensure known value
	if monitor.GeneralRegion != "" {
		data.GeneralRegion = types.StringValue(monitor.GeneralRegion)
	} else if data.GeneralRegion.IsUnknown() {
		data.GeneralRegion = types.StringNull()
	}
	if monitor.SpecificRegion != "" {
		data.SpecificRegion = types.StringValue(monitor.SpecificRegion)
	} else if data.SpecificRegion.IsUnknown() {
		data.SpecificRegion = types.StringNull()
	}
	// Computed field - must always be set to a known value
	if monitor.LastChecked != "" {
		data.LastChecked = types.StringValue(normalizeTimestamp(monitor.LastChecked))
	} else {
		data.LastChecked = types.StringNull()
	}

	// HTTP specific
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

	// DNS specific
	if monitor.DNSRecordType != "" {
		data.DNSRecordType = types.StringValue(monitor.DNSRecordType)
	}
	if monitor.ExpectedValue != "" {
		data.ExpectedValue = types.StringValue(monitor.ExpectedValue)
	}
	if monitor.Nameserver != "" {
		data.Nameserver = types.StringValue(monitor.Nameserver)
	}

	// TCP specific
	if monitor.Host != "" {
		data.Host = types.StringValue(monitor.Host)
	}
	if monitor.Port != 0 {
		data.Port = types.Int64Value(int64(monitor.Port))
	}

	// SSL specific
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
}

// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

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
var _ resource.Resource = &AlertResource{}
var _ resource.ResourceWithImportState = &AlertResource{}

func NewAlertResource() resource.Resource {
	return &AlertResource{}
}

// AlertResource defines the resource implementation.
type AlertResource struct {
	client *client.Client
}

// AlertResourceModel describes the resource data model.
type AlertResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	MonitorID          types.String `tfsdk:"monitor_id"`
	Type               types.String `tfsdk:"type"`
	Target             types.String `tfsdk:"target"`
	IsEnabled          types.Bool   `tfsdk:"is_enabled"`
	TriggerThreshold   types.Int64  `tfsdk:"trigger_threshold"`
	RecoveryThreshold  types.Int64  `tfsdk:"recovery_threshold"`
	MinIntervalMinutes types.Int64  `tfsdk:"min_interval_minutes"`
	CustomMessage      types.String `tfsdk:"custom_message"`
	IncludeDetails     types.Bool   `tfsdk:"include_details"`
	LastTriggeredAt    types.String `tfsdk:"last_triggered_at"`
	CreatedAt          types.String `tfsdk:"created_at"`
	UpdatedAt          types.String `tfsdk:"updated_at"`
}

func (r *AlertResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert"
}

func (r *AlertResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an alert configuration for a monitor on ackack.io.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the alert.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"monitor_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the monitor this alert is attached to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of alert. Must be one of: `email`, `webhook`, `discord`, `slack`, `pagerduty`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("email", "webhook", "discord", "slack", "pagerduty"),
				},
			},
			"target": schema.StringAttribute{
				MarkdownDescription: "The target for the alert (email address, webhook URL, etc.).",
				Required:            true,
			},
			"is_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the alert is enabled. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"trigger_threshold": schema.Int64Attribute{
				MarkdownDescription: "Number of consecutive failures before triggering the alert. Defaults to `1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"recovery_threshold": schema.Int64Attribute{
				MarkdownDescription: "Number of consecutive successes before sending recovery notification. Defaults to `1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"min_interval_minutes": schema.Int64Attribute{
				MarkdownDescription: "Minimum interval between alerts, in minutes. Defaults to `5`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(5),
			},
			"custom_message": schema.StringAttribute{
				MarkdownDescription: "Custom message to include in alerts.",
				Optional:            true,
			},
			"include_details": schema.BoolAttribute{
				MarkdownDescription: "Whether to include detailed information in the alert.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"last_triggered_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the alert was last triggered.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the alert was created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the alert was last updated.",
				Computed:            true,
			},
		},
	}
}

func (r *AlertResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AlertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AlertResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.CreateAlertRequest{
		MonitorID: data.MonitorID.ValueString(),
		Type:      data.Type.ValueString(),
		Target:    data.Target.ValueString(),
	}

	if !data.IsEnabled.IsNull() {
		isEnabled := data.IsEnabled.ValueBool()
		createReq.IsEnabled = &isEnabled
	}
	if !data.TriggerThreshold.IsNull() {
		createReq.TriggerThreshold = int(data.TriggerThreshold.ValueInt64())
	}
	if !data.RecoveryThreshold.IsNull() {
		createReq.RecoveryThreshold = int(data.RecoveryThreshold.ValueInt64())
	}
	if !data.MinIntervalMinutes.IsNull() {
		createReq.MinIntervalMinutes = int(data.MinIntervalMinutes.ValueInt64())
	}
	if !data.CustomMessage.IsNull() {
		createReq.CustomMessage = data.CustomMessage.ValueString()
	}
	if !data.IncludeDetails.IsNull() {
		includeDetails := data.IncludeDetails.ValueBool()
		createReq.IncludeDetails = &includeDetails
	}

	alert, err := r.client.CreateAlert(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create alert, got error: %s", err))
		return
	}

	r.updateModelFromResponse(&data, alert)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AlertResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alert, err := r.client.GetAlert(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert, got error: %s", err))
		return
	}

	r.updateModelFromResponse(&data, alert)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AlertResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := client.UpdateAlertRequest{
		Target: data.Target.ValueString(),
	}

	if !data.IsEnabled.IsNull() {
		isEnabled := data.IsEnabled.ValueBool()
		updateReq.IsEnabled = &isEnabled
	}
	if !data.TriggerThreshold.IsNull() {
		updateReq.TriggerThreshold = int(data.TriggerThreshold.ValueInt64())
	}
	if !data.RecoveryThreshold.IsNull() {
		updateReq.RecoveryThreshold = int(data.RecoveryThreshold.ValueInt64())
	}
	if !data.MinIntervalMinutes.IsNull() {
		updateReq.MinIntervalMinutes = int(data.MinIntervalMinutes.ValueInt64())
	}
	if !data.CustomMessage.IsNull() {
		updateReq.CustomMessage = data.CustomMessage.ValueString()
	}
	if !data.IncludeDetails.IsNull() {
		includeDetails := data.IncludeDetails.ValueBool()
		updateReq.IncludeDetails = &includeDetails
	}

	alert, err := r.client.UpdateAlert(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update alert, got error: %s", err))
		return
	}

	r.updateModelFromResponse(&data, alert)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AlertResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAlert(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete alert, got error: %s", err))
		return
	}
}

func (r *AlertResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AlertResource) updateModelFromResponse(data *AlertResourceModel, alert *client.Alert) {
	data.ID = types.StringValue(alert.ID)
	data.MonitorID = types.StringValue(alert.MonitorID)
	data.Type = types.StringValue(alert.Type)
	data.Target = types.StringValue(alert.Target)
	data.IsEnabled = types.BoolValue(alert.IsEnabled)
	data.TriggerThreshold = types.Int64Value(int64(alert.TriggerThreshold))
	data.RecoveryThreshold = types.Int64Value(int64(alert.RecoveryThreshold))
	data.MinIntervalMinutes = types.Int64Value(int64(alert.MinIntervalMinutes))
	data.IncludeDetails = types.BoolValue(alert.IncludeDetails)
	data.CreatedAt = types.StringValue(alert.CreatedAt)
	data.UpdatedAt = types.StringValue(alert.UpdatedAt)

	if alert.CustomMessage != "" {
		data.CustomMessage = types.StringValue(alert.CustomMessage)
	}
	if alert.LastTriggeredAt != "" {
		data.LastTriggeredAt = types.StringValue(alert.LastTriggeredAt)
	}
}

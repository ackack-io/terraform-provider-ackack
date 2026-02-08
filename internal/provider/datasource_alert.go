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
var _ datasource.DataSource = &AlertDataSource{}

func NewAlertDataSource() datasource.DataSource {
	return &AlertDataSource{}
}

// AlertDataSource defines the data source implementation.
type AlertDataSource struct {
	client *client.Client
}

// AlertDataSourceModel describes the data source data model.
type AlertDataSourceModel struct {
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

func (d *AlertDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert"
}

func (d *AlertDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a specific alert.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the alert.",
				Required:            true,
			},
			"monitor_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the monitor this alert is attached to.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of alert (email, webhook, discord, slack, pagerduty).",
				Computed:            true,
			},
			"target": schema.StringAttribute{
				MarkdownDescription: "The target for the alert.",
				Computed:            true,
			},
			"is_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the alert is enabled.",
				Computed:            true,
			},
			"trigger_threshold": schema.Int64Attribute{
				MarkdownDescription: "Number of consecutive failures before triggering the alert.",
				Computed:            true,
			},
			"recovery_threshold": schema.Int64Attribute{
				MarkdownDescription: "Number of consecutive successes before sending recovery notification.",
				Computed:            true,
			},
			"min_interval_minutes": schema.Int64Attribute{
				MarkdownDescription: "Minimum interval between alerts, in minutes.",
				Computed:            true,
			},
			"custom_message": schema.StringAttribute{
				MarkdownDescription: "Custom message to include in alerts.",
				Computed:            true,
			},
			"include_details": schema.BoolAttribute{
				MarkdownDescription: "Whether to include detailed information in the alert.",
				Computed:            true,
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

func (d *AlertDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AlertDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alert, err := d.client.GetAlert(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alert, got error: %s", err))
		return
	}

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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

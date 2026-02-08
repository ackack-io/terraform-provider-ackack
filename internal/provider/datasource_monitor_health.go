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
var _ datasource.DataSource = &MonitorHealthDataSource{}

func NewMonitorHealthDataSource() datasource.DataSource {
	return &MonitorHealthDataSource{}
}

// MonitorHealthDataSource defines the data source implementation.
type MonitorHealthDataSource struct {
	client *client.Client
}

// MonitorHealthDataSourceModel describes the data source data model.
type MonitorHealthDataSourceModel struct {
	MonitorID       types.String  `tfsdk:"monitor_id"`
	MonitorName     types.String  `tfsdk:"monitor_name"`
	IsInFlight      types.Bool    `tfsdk:"is_in_flight"`
	InFlightSeconds types.Float64 `tfsdk:"in_flight_seconds"`
	Throttled       types.Bool    `tfsdk:"throttled"`
	ThrottleReason  types.String  `tfsdk:"throttle_reason"`
	DampeningLevel  types.Int64   `tfsdk:"dampening_level"`
	DampeningName   types.String  `tfsdk:"dampening_name"`
	DampeningReason types.String  `tfsdk:"dampening_reason"`
	FailureRate     types.Float64 `tfsdk:"failure_rate"`
	P95LatencyMs    types.Int64   `tfsdk:"p95_latency_ms"`
	StuckCount      types.Int64   `tfsdk:"stuck_count"`
}

func (d *MonitorHealthDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor_health"
}

func (d *MonitorHealthDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get health information for a specific monitor.",

		Attributes: map[string]schema.Attribute{
			"monitor_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the monitor.",
				Required:            true,
			},
			"monitor_name": schema.StringAttribute{
				MarkdownDescription: "The name of the monitor.",
				Computed:            true,
			},
			"is_in_flight": schema.BoolAttribute{
				MarkdownDescription: "Whether a check is currently in progress.",
				Computed:            true,
			},
			"in_flight_seconds": schema.Float64Attribute{
				MarkdownDescription: "How long the current check has been running.",
				Computed:            true,
			},
			"throttled": schema.BoolAttribute{
				MarkdownDescription: "Whether the monitor is being throttled.",
				Computed:            true,
			},
			"throttle_reason": schema.StringAttribute{
				MarkdownDescription: "The reason for throttling.",
				Computed:            true,
			},
			"dampening_level": schema.Int64Attribute{
				MarkdownDescription: "The current dampening level.",
				Computed:            true,
			},
			"dampening_name": schema.StringAttribute{
				MarkdownDescription: "The name of the dampening state.",
				Computed:            true,
			},
			"dampening_reason": schema.StringAttribute{
				MarkdownDescription: "The reason for dampening.",
				Computed:            true,
			},
			"failure_rate": schema.Float64Attribute{
				MarkdownDescription: "The recent failure rate.",
				Computed:            true,
			},
			"p95_latency_ms": schema.Int64Attribute{
				MarkdownDescription: "The 95th percentile latency in milliseconds.",
				Computed:            true,
			},
			"stuck_count": schema.Int64Attribute{
				MarkdownDescription: "The number of stuck checks.",
				Computed:            true,
			},
		},
	}
}

func (d *MonitorHealthDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonitorHealthDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MonitorHealthDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	health, err := d.client.GetMonitorHealth(ctx, data.MonitorID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get monitor health, got error: %s", err))
		return
	}

	data.MonitorName = types.StringValue(health.MonitorName)
	data.IsInFlight = types.BoolValue(health.IsInFlight)
	data.InFlightSeconds = types.Float64Value(health.InFlightSeconds)
	data.Throttled = types.BoolValue(health.Throttled)
	data.DampeningLevel = types.Int64Value(int64(health.DampeningLevel))
	data.FailureRate = types.Float64Value(health.FailureRate)
	data.P95LatencyMs = types.Int64Value(int64(health.P95LatencyMs))
	data.StuckCount = types.Int64Value(int64(health.StuckCount))

	if health.ThrottleReason != "" {
		data.ThrottleReason = types.StringValue(health.ThrottleReason)
	}
	if health.DampeningName != "" {
		data.DampeningName = types.StringValue(health.DampeningName)
	}
	if health.DampeningReason != "" {
		data.DampeningReason = types.StringValue(health.DampeningReason)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

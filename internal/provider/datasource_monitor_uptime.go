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
var _ datasource.DataSource = &MonitorUptimeDataSource{}

func NewMonitorUptimeDataSource() datasource.DataSource {
	return &MonitorUptimeDataSource{}
}

// MonitorUptimeDataSource defines the data source implementation.
type MonitorUptimeDataSource struct {
	client *client.Client
}

// MonitorUptimeDataSourceModel describes the data source data model.
type MonitorUptimeDataSourceModel struct {
	MonitorID types.String  `tfsdk:"monitor_id"`
	Hours     types.Int64   `tfsdk:"hours"`
	Uptime    types.Float64 `tfsdk:"uptime"`
}

func (d *MonitorUptimeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor_uptime"
}

func (d *MonitorUptimeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get the uptime percentage for a monitor over a specified time window.",

		Attributes: map[string]schema.Attribute{
			"monitor_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the monitor.",
				Required:            true,
			},
			"hours": schema.Int64Attribute{
				MarkdownDescription: "The time window in hours. Default is 24.",
				Optional:            true,
			},
			"uptime": schema.Float64Attribute{
				MarkdownDescription: "The uptime percentage.",
				Computed:            true,
			},
		},
	}
}

func (d *MonitorUptimeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonitorUptimeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MonitorUptimeDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hours := 0
	if !data.Hours.IsNull() {
		hours = int(data.Hours.ValueInt64())
	}

	uptimeResp, err := d.client.GetMonitorUptime(ctx, data.MonitorID.ValueString(), hours)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get monitor uptime, got error: %s", err))
		return
	}

	data.Hours = types.Int64Value(int64(uptimeResp.Hours))
	data.Uptime = types.Float64Value(uptimeResp.Uptime)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

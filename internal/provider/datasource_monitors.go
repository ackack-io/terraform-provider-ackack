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
var _ datasource.DataSource = &MonitorsDataSource{}

func NewMonitorsDataSource() datasource.DataSource {
	return &MonitorsDataSource{}
}

// MonitorsDataSource defines the data source implementation.
type MonitorsDataSource struct {
	client *client.Client
}

// MonitorsDataSourceModel describes the data source data model.
type MonitorsDataSourceModel struct {
	Monitors []MonitorListItemModel `tfsdk:"monitors"`
}

// MonitorListItemModel describes a single monitor in the list.
type MonitorListItemModel struct {
	ID               types.String  `tfsdk:"id"`
	Name             types.String  `tfsdk:"name"`
	Type             types.String  `tfsdk:"type"`
	IsEnabled        types.Bool    `tfsdk:"is_enabled"`
	Status           types.String  `tfsdk:"status"`
	UptimePercentage types.Float64 `tfsdk:"uptime_percentage"`
	LastChecked      types.String  `tfsdk:"last_checked"`
	CreatedAt        types.String  `tfsdk:"created_at"`
}

func (d *MonitorsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitors"
}

func (d *MonitorsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to list all monitors.",

		Attributes: map[string]schema.Attribute{
			"monitors": schema.ListNestedAttribute{
				MarkdownDescription: "List of monitors.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the monitor.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the monitor.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of monitor.",
							Computed:            true,
						},
						"is_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether the monitor is enabled.",
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
					},
				},
			},
		},
	}
}

func (d *MonitorsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonitorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MonitorsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	monitors, err := d.client.ListMonitors(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list monitors, got error: %s", err))
		return
	}

	data.Monitors = make([]MonitorListItemModel, len(monitors))
	for i, monitor := range monitors {
		data.Monitors[i] = MonitorListItemModel{
			ID:               types.StringValue(monitor.ID),
			Name:             types.StringValue(monitor.Name),
			Type:             types.StringValue(monitor.Type),
			IsEnabled:        types.BoolValue(monitor.IsEnabled),
			Status:           types.StringValue(monitor.Status),
			UptimePercentage: types.Float64Value(monitor.UptimePercentage),
			CreatedAt:        types.StringValue(monitor.CreatedAt),
		}
		if monitor.LastChecked != "" {
			data.Monitors[i].LastChecked = types.StringValue(monitor.LastChecked)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

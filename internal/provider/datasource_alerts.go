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
var _ datasource.DataSource = &AlertsDataSource{}

func NewAlertsDataSource() datasource.DataSource {
	return &AlertsDataSource{}
}

// AlertsDataSource defines the data source implementation.
type AlertsDataSource struct {
	client *client.Client
}

// AlertsDataSourceModel describes the data source data model.
type AlertsDataSourceModel struct {
	Alerts []AlertListItemModel `tfsdk:"alerts"`
}

// AlertListItemModel describes a single alert in the list.
type AlertListItemModel struct {
	ID              types.String `tfsdk:"id"`
	MonitorID       types.String `tfsdk:"monitor_id"`
	Type            types.String `tfsdk:"type"`
	Target          types.String `tfsdk:"target"`
	IsEnabled       types.Bool   `tfsdk:"is_enabled"`
	LastTriggeredAt types.String `tfsdk:"last_triggered_at"`
	CreatedAt       types.String `tfsdk:"created_at"`
}

func (d *AlertsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alerts"
}

func (d *AlertsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to list all alerts.",

		Attributes: map[string]schema.Attribute{
			"alerts": schema.ListNestedAttribute{
				MarkdownDescription: "List of alerts.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the alert.",
							Computed:            true,
						},
						"monitor_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the monitor this alert is attached to.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of alert.",
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
						"last_triggered_at": schema.StringAttribute{
							MarkdownDescription: "The timestamp when the alert was last triggered.",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "The timestamp when the alert was created.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *AlertsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AlertsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AlertsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alerts, err := d.client.ListAlerts(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list alerts, got error: %s", err))
		return
	}

	data.Alerts = make([]AlertListItemModel, len(alerts))
	for i, alert := range alerts {
		data.Alerts[i] = AlertListItemModel{
			ID:        types.StringValue(alert.ID),
			MonitorID: types.StringValue(alert.MonitorID),
			Type:      types.StringValue(alert.Type),
			Target:    types.StringValue(alert.Target),
			IsEnabled: types.BoolValue(alert.IsEnabled),
			CreatedAt: types.StringValue(alert.CreatedAt),
		}
		if alert.LastTriggeredAt != "" {
			data.Alerts[i].LastTriggeredAt = types.StringValue(alert.LastTriggeredAt)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

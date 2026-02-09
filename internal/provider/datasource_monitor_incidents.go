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
var _ datasource.DataSource = &MonitorIncidentsDataSource{}

func NewMonitorIncidentsDataSource() datasource.DataSource {
	return &MonitorIncidentsDataSource{}
}

// MonitorIncidentsDataSource defines the data source implementation.
type MonitorIncidentsDataSource struct {
	client *client.Client
}

// MonitorIncidentsDataSourceModel describes the data source data model.
type MonitorIncidentsDataSourceModel struct {
	MonitorID types.String        `tfsdk:"monitor_id"`
	Limit     types.Int64         `tfsdk:"limit"`
	Incidents []IncidentItemModel `tfsdk:"incidents"`
}

// IncidentItemModel describes a single incident.
type IncidentItemModel struct {
	ID              types.String `tfsdk:"id"`
	Status          types.String `tfsdk:"status"`
	Severity        types.String `tfsdk:"severity"`
	Summary         types.String `tfsdk:"summary"`
	Details         types.String `tfsdk:"details"`
	StartedAt       types.String `tfsdk:"started_at"`
	ResolvedAt      types.String `tfsdk:"resolved_at"`
	DurationSeconds types.Int64  `tfsdk:"duration_seconds"`
	Notified        types.Bool   `tfsdk:"notified"`
}

func (d *MonitorIncidentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor_incidents"
}

func (d *MonitorIncidentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get incidents for a monitor.",

		Attributes: map[string]schema.Attribute{
			"monitor_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the monitor.",
				Required:            true,
			},
			"limit": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of incidents to return. Default is 50, max is 500.",
				Optional:            true,
			},
			"incidents": schema.ListNestedAttribute{
				MarkdownDescription: "List of incidents.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The incident ID.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "The incident status.",
							Computed:            true,
						},
						"severity": schema.StringAttribute{
							MarkdownDescription: "The incident severity.",
							Computed:            true,
						},
						"summary": schema.StringAttribute{
							MarkdownDescription: "A summary of the incident.",
							Computed:            true,
						},
						"details": schema.StringAttribute{
							MarkdownDescription: "Details about the incident.",
							Computed:            true,
						},
						"started_at": schema.StringAttribute{
							MarkdownDescription: "When the incident started.",
							Computed:            true,
						},
						"resolved_at": schema.StringAttribute{
							MarkdownDescription: "When the incident was resolved.",
							Computed:            true,
						},
						"duration_seconds": schema.Int64Attribute{
							MarkdownDescription: "Duration of the incident in seconds.",
							Computed:            true,
						},
						"notified": schema.BoolAttribute{
							MarkdownDescription: "Whether notifications were sent.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *MonitorIncidentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonitorIncidentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MonitorIncidentsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	limit := 0
	if !data.Limit.IsNull() {
		limit = int(data.Limit.ValueInt64())
	}

	incidents, err := d.client.GetMonitorIncidents(ctx, data.MonitorID.ValueString(), limit)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get monitor incidents, got error: %s", err))
		return
	}

	data.Incidents = make([]IncidentItemModel, len(incidents))
	for i, incident := range incidents {
		data.Incidents[i] = IncidentItemModel{
			ID:              types.StringValue(incident.ID),
			Status:          types.StringValue(incident.Status),
			Severity:        types.StringValue(incident.Severity),
			StartedAt:       types.StringValue(incident.StartedAt),
			DurationSeconds: types.Int64Value(int64(incident.DurationSeconds)),
			Notified:        types.BoolValue(incident.Notified),
		}
		if incident.Summary != "" {
			data.Incidents[i].Summary = types.StringValue(incident.Summary)
		}
		if incident.Details != "" {
			data.Incidents[i].Details = types.StringValue(incident.Details)
		}
		if incident.ResolvedAt != "" {
			data.Incidents[i].ResolvedAt = types.StringValue(incident.ResolvedAt)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

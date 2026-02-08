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
var _ datasource.DataSource = &SystemsDataSource{}

func NewSystemsDataSource() datasource.DataSource {
	return &SystemsDataSource{}
}

// SystemsDataSource defines the data source implementation.
type SystemsDataSource struct {
	client *client.Client
}

// SystemsDataSourceModel describes the data source data model.
type SystemsDataSourceModel struct {
	Systems []SystemListItemModel `tfsdk:"systems"`
}

// SystemListItemModel describes a single system in the list.
type SystemListItemModel struct {
	ID            types.String  `tfsdk:"id"`
	Name          types.String  `tfsdk:"name"`
	Status        types.String  `tfsdk:"status"`
	MonitorCount  types.Int64   `tfsdk:"monitor_count"`
	HealthyCount  types.Int64   `tfsdk:"healthy_count"`
	OverallUptime types.Float64 `tfsdk:"overall_uptime"`
	CreatedAt     types.String  `tfsdk:"created_at"`
}

func (d *SystemsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_systems"
}

func (d *SystemsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to list all systems.",

		Attributes: map[string]schema.Attribute{
			"systems": schema.ListNestedAttribute{
				MarkdownDescription: "List of systems.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the system.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the system.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "The current status of the system.",
							Computed:            true,
						},
						"monitor_count": schema.Int64Attribute{
							MarkdownDescription: "The number of monitors in the system.",
							Computed:            true,
						},
						"healthy_count": schema.Int64Attribute{
							MarkdownDescription: "The number of healthy monitors in the system.",
							Computed:            true,
						},
						"overall_uptime": schema.Float64Attribute{
							MarkdownDescription: "The overall uptime percentage of the system.",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "The timestamp when the system was created.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *SystemsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SystemsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	systems, err := d.client.ListSystems(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list systems, got error: %s", err))
		return
	}

	data.Systems = make([]SystemListItemModel, len(systems))
	for i, system := range systems {
		data.Systems[i] = SystemListItemModel{
			ID:            types.StringValue(system.ID),
			Name:          types.StringValue(system.Name),
			Status:        types.StringValue(system.Status),
			MonitorCount:  types.Int64Value(int64(system.MonitorCount)),
			HealthyCount:  types.Int64Value(int64(system.HealthyCount)),
			OverallUptime: types.Float64Value(system.OverallUptime),
			CreatedAt:     types.StringValue(system.CreatedAt),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

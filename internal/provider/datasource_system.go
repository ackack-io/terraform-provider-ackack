// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/ackack-io/terraform-provider-ackack/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SystemDataSource{}

func NewSystemDataSource() datasource.DataSource {
	return &SystemDataSource{}
}

// SystemDataSource defines the data source implementation.
type SystemDataSource struct {
	client *client.Client
}

// SystemDataSourceModel describes the data source data model.
type SystemDataSourceModel struct {
	ID            types.String  `tfsdk:"id"`
	Name          types.String  `tfsdk:"name"`
	Description   types.String  `tfsdk:"description"`
	Priority      types.String  `tfsdk:"priority"`
	Status        types.String  `tfsdk:"status"`
	ExternalLinks types.List    `tfsdk:"external_links"`
	MonitorCount  types.Int64   `tfsdk:"monitor_count"`
	HealthyCount  types.Int64   `tfsdk:"healthy_count"`
	DegradedCount types.Int64   `tfsdk:"degraded_count"`
	ErrorCount    types.Int64   `tfsdk:"error_count"`
	OverallUptime types.Float64 `tfsdk:"overall_uptime"`
	CreatedAt     types.String  `tfsdk:"created_at"`
	UpdatedAt     types.String  `tfsdk:"updated_at"`
}

func (d *SystemDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system"
}

func (d *SystemDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a specific system.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the system.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the system.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the system.",
				Computed:            true,
			},
			"priority": schema.StringAttribute{
				MarkdownDescription: "The priority of the system.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The current status of the system.",
				Computed:            true,
			},
			"external_links": schema.ListNestedAttribute{
				MarkdownDescription: "External links associated with this system.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the link.",
							Computed:            true,
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "The URL of the link.",
							Computed:            true,
						},
					},
				},
			},
			"monitor_count": schema.Int64Attribute{
				MarkdownDescription: "The number of monitors in the system.",
				Computed:            true,
			},
			"healthy_count": schema.Int64Attribute{
				MarkdownDescription: "The number of healthy monitors in the system.",
				Computed:            true,
			},
			"degraded_count": schema.Int64Attribute{
				MarkdownDescription: "The number of degraded monitors in the system.",
				Computed:            true,
			},
			"error_count": schema.Int64Attribute{
				MarkdownDescription: "The number of monitors in error state.",
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
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the system was last updated.",
				Computed:            true,
			},
		},
	}
}

func (d *SystemDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SystemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	system, err := d.client.GetSystem(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read system, got error: %s", err))
		return
	}

	data.Name = types.StringValue(system.Name)
	data.Status = types.StringValue(system.Status)
	data.MonitorCount = types.Int64Value(int64(system.MonitorCount))
	data.HealthyCount = types.Int64Value(int64(system.HealthyCount))
	data.DegradedCount = types.Int64Value(int64(system.DegradedCount))
	data.ErrorCount = types.Int64Value(int64(system.ErrorCount))
	data.OverallUptime = types.Float64Value(system.OverallUptime)
	data.CreatedAt = types.StringValue(system.CreatedAt)
	data.UpdatedAt = types.StringValue(system.UpdatedAt)

	if system.Description != "" {
		data.Description = types.StringValue(system.Description)
	}
	if system.Priority != "" {
		data.Priority = types.StringValue(system.Priority)
	}

	// Convert external links
	if len(system.ExternalLinks) > 0 {
		linkObjects := make([]attr.Value, len(system.ExternalLinks))
		for i, link := range system.ExternalLinks {
			linkObj, diags := types.ObjectValue(
				map[string]attr.Type{
					"name": types.StringType,
					"url":  types.StringType,
				},
				map[string]attr.Value{
					"name": types.StringValue(link.Name),
					"url":  types.StringValue(link.URL),
				},
			)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			linkObjects[i] = linkObj
		}
		linksList, diags := types.ListValue(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name": types.StringType,
				"url":  types.StringType,
			},
		}, linkObjects)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.ExternalLinks = linksList
	} else {
		data.ExternalLinks = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name": types.StringType,
				"url":  types.StringType,
			},
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

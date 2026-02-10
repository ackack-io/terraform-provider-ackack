// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/ackack-io/terraform-provider-ackack/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SystemResource{}
var _ resource.ResourceWithImportState = &SystemResource{}

func NewSystemResource() resource.Resource {
	return &SystemResource{}
}

// SystemResource defines the resource implementation.
type SystemResource struct {
	client *client.Client
}

// SystemResourceModel describes the resource data model.
type SystemResourceModel struct {
	ID            types.String  `tfsdk:"id"`
	Name          types.String  `tfsdk:"name"`
	Description   types.String  `tfsdk:"description"`
	Priority      types.String  `tfsdk:"priority"`
	Status        types.String  `tfsdk:"status"`
	MonitorIDs    types.Set     `tfsdk:"monitor_ids"`
	ExternalLinks types.List    `tfsdk:"external_links"`
	MonitorCount  types.Int64   `tfsdk:"monitor_count"`
	HealthyCount  types.Int64   `tfsdk:"healthy_count"`
	OverallUptime types.Float64 `tfsdk:"overall_uptime"`
	CreatedAt     types.String  `tfsdk:"created_at"`
	UpdatedAt     types.String  `tfsdk:"updated_at"`
}

// ExternalLinkModel describes an external link.
type ExternalLinkModel struct {
	Name types.String `tfsdk:"name"`
	URL  types.String `tfsdk:"url"`
}

func (r *SystemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system"
}

func (r *SystemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a system grouping of monitors on ackack.io.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the system.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the system.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the system.",
				Optional:            true,
			},
			"priority": schema.StringAttribute{
				MarkdownDescription: "The priority of the system.",
				Optional:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The current status of the system.",
				Computed:            true,
			},
			"monitor_ids": schema.SetAttribute{
				MarkdownDescription: "The IDs of monitors in this system. At least one monitor is required.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"external_links": schema.ListNestedAttribute{
				MarkdownDescription: "External links associated with this system.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the link.",
							Required:            true,
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "The URL of the link.",
							Required:            true,
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

func (r *SystemResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SystemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SystemResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract monitor IDs
	var monitorIDs []string
	resp.Diagnostics.Append(data.MonitorIDs.ElementsAs(ctx, &monitorIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external links
	var externalLinks []client.ExternalLink
	if !data.ExternalLinks.IsNull() {
		var linkModels []ExternalLinkModel
		resp.Diagnostics.Append(data.ExternalLinks.ElementsAs(ctx, &linkModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, lm := range linkModels {
			externalLinks = append(externalLinks, client.ExternalLink{
				Name: lm.Name.ValueString(),
				URL:  lm.URL.ValueString(),
			})
		}
	}

	createReq := client.CreateSystemRequest{
		Name:          data.Name.ValueString(),
		MonitorIDs:    monitorIDs,
		ExternalLinks: externalLinks,
	}

	if !data.Description.IsNull() {
		createReq.Description = data.Description.ValueString()
	}
	if !data.Priority.IsNull() {
		createReq.Priority = data.Priority.ValueString()
	}

	system, err := r.client.CreateSystem(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create system, got error: %s", err))
		return
	}

	// After creating, fetch the system with stats
	systemWithStats, err := r.client.GetSystem(ctx, system.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read created system, got error: %s", err))
		return
	}

	r.updateModelFromResponse(ctx, &data, systemWithStats, monitorIDs)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SystemResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current monitor IDs from state for comparison
	var currentMonitorIDs []string
	if !data.MonitorIDs.IsNull() {
		resp.Diagnostics.Append(data.MonitorIDs.ElementsAs(ctx, &currentMonitorIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	system, err := r.client.GetSystem(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read system, got error: %s", err))
		return
	}

	r.updateModelFromResponse(ctx, &data, system, currentMonitorIDs)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SystemResourceModel
	var state SystemResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract new monitor IDs
	var newMonitorIDs []string
	resp.Diagnostics.Append(data.MonitorIDs.ElementsAs(ctx, &newMonitorIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract old monitor IDs
	var oldMonitorIDs []string
	resp.Diagnostics.Append(state.MonitorIDs.ElementsAs(ctx, &oldMonitorIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external links
	var externalLinks []client.ExternalLink
	if !data.ExternalLinks.IsNull() {
		var linkModels []ExternalLinkModel
		resp.Diagnostics.Append(data.ExternalLinks.ElementsAs(ctx, &linkModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, lm := range linkModels {
			externalLinks = append(externalLinks, client.ExternalLink{
				Name: lm.Name.ValueString(),
				URL:  lm.URL.ValueString(),
			})
		}
	}

	// Update system metadata
	updateReq := client.UpdateSystemRequest{
		Name:          data.Name.ValueString(),
		ExternalLinks: externalLinks,
	}

	if !data.Description.IsNull() {
		updateReq.Description = data.Description.ValueString()
	}
	if !data.Priority.IsNull() {
		updateReq.Priority = data.Priority.ValueString()
	}

	_, err := r.client.UpdateSystem(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update system, got error: %s", err))
		return
	}

	// Calculate monitor changes
	toAdd := difference(newMonitorIDs, oldMonitorIDs)
	toRemove := difference(oldMonitorIDs, newMonitorIDs)

	// Add new monitors
	if len(toAdd) > 0 {
		err = r.client.AddMonitorsToSystem(ctx, data.ID.ValueString(), toAdd)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add monitors to system, got error: %s", err))
			return
		}
	}

	// Remove old monitors
	if len(toRemove) > 0 {
		err = r.client.RemoveMonitorsFromSystem(ctx, data.ID.ValueString(), toRemove)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove monitors from system, got error: %s", err))
			return
		}
	}

	// Fetch updated system
	system, err := r.client.GetSystem(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read updated system, got error: %s", err))
		return
	}

	r.updateModelFromResponse(ctx, &data, system, newMonitorIDs)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SystemResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSystem(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete system, got error: %s", err))
		return
	}
}

func (r *SystemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *SystemResource) updateModelFromResponse(ctx context.Context, data *SystemResourceModel, system *client.SystemWithStats, monitorIDs []string) {
	data.ID = types.StringValue(system.ID)
	data.Name = types.StringValue(system.Name)
	data.Status = types.StringValue(system.Status)
	data.MonitorCount = types.Int64Value(int64(system.MonitorCount))
	data.HealthyCount = types.Int64Value(int64(system.HealthyCount))
	data.OverallUptime = types.Float64Value(system.OverallUptime)
	data.CreatedAt = types.StringValue(system.CreatedAt)
	data.UpdatedAt = types.StringValue(system.UpdatedAt)

	if system.Description != "" {
		data.Description = types.StringValue(system.Description)
	}
	if system.Priority != "" {
		data.Priority = types.StringValue(system.Priority)
	}

	// Preserve the monitor_ids from the plan/state since the API doesn't return them
	if len(monitorIDs) > 0 {
		monitorIDsSet, d := types.SetValueFrom(ctx, types.StringType, monitorIDs)
		if d.HasError() {
			return
		}
		data.MonitorIDs = monitorIDsSet
	}

	// Convert external links
	if len(system.ExternalLinks) > 0 {
		linkObjects := make([]attr.Value, len(system.ExternalLinks))
		for i, link := range system.ExternalLinks {
			linkObj, d := types.ObjectValue(
				map[string]attr.Type{
					"name": types.StringType,
					"url":  types.StringType,
				},
				map[string]attr.Value{
					"name": types.StringValue(link.Name),
					"url":  types.StringValue(link.URL),
				},
			)
			if d.HasError() {
				return
			}
			linkObjects[i] = linkObj
		}
		linksList, d := types.ListValue(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name": types.StringType,
				"url":  types.StringType,
			},
		}, linkObjects)
		if d.HasError() {
			return
		}
		data.ExternalLinks = linksList
	}
}

// difference returns elements in a that are not in b.
func difference(a, b []string) []string {
	bMap := make(map[string]bool)
	for _, v := range b {
		bMap[v] = true
	}
	var result []string
	for _, v := range a {
		if !bMap[v] {
			result = append(result, v)
		}
	}
	return result
}

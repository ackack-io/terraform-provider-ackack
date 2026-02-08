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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ReportResource{}
var _ resource.ResourceWithImportState = &ReportResource{}

func NewReportResource() resource.Resource {
	return &ReportResource{}
}

// ReportResource defines the resource implementation.
type ReportResource struct {
	client *client.Client
}

// ReportResourceModel describes the resource data model.
type ReportResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	ReportType  types.String `tfsdk:"report_type"`
	Format      types.String `tfsdk:"format"`
	StartTime   types.String `tfsdk:"start_time"`
	EndTime     types.String `tfsdk:"end_time"`
	MonitorIDs  types.Set    `tfsdk:"monitor_ids"`
	SystemIDs   types.Set    `tfsdk:"system_ids"`
	Metrics     types.String `tfsdk:"metrics"`
	Status      types.String `tfsdk:"status"`
	FilePath    types.String `tfsdk:"file_path"`
	CompletedAt types.String `tfsdk:"completed_at"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func (r *ReportResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_report"
}

func (r *ReportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a report on ackack.io. Reports cannot be updated - any configuration change will trigger replacement.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the report.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the report.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"report_type": schema.StringAttribute{
				MarkdownDescription: "The type of report. Must be one of: `uptime`, `incidents`, `custom`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("uptime", "incidents", "custom"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"format": schema.StringAttribute{
				MarkdownDescription: "The format of the report. Must be one of: `pdf`, `csv`, `json`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("pdf", "csv", "json"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"start_time": schema.StringAttribute{
				MarkdownDescription: "The start time for the report in ISO 8601 format.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"end_time": schema.StringAttribute{
				MarkdownDescription: "The end time for the report in ISO 8601 format.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"monitor_ids": schema.SetAttribute{
				MarkdownDescription: "The IDs of monitors to include in the report. If not specified, all monitors are included.",
				Optional:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setRequiresReplace(),
				},
			},
			"system_ids": schema.SetAttribute{
				MarkdownDescription: "The IDs of systems to include in the report.",
				Optional:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setRequiresReplace(),
				},
			},
			"metrics": schema.StringAttribute{
				MarkdownDescription: "Custom metrics configuration as a JSON string.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the report.",
				Computed:            true,
			},
			"file_path": schema.StringAttribute{
				MarkdownDescription: "The path to the generated report file.",
				Computed:            true,
			},
			"completed_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the report was completed.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the report was created.",
				Computed:            true,
			},
		},
	}
}

func (r *ReportResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ReportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ReportResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.CreateReportRequest{
		Name:       data.Name.ValueString(),
		ReportType: data.ReportType.ValueString(),
		Format:     data.Format.ValueString(),
		StartTime:  data.StartTime.ValueString(),
		EndTime:    data.EndTime.ValueString(),
	}

	if !data.MonitorIDs.IsNull() {
		var monitorIDs []string
		resp.Diagnostics.Append(data.MonitorIDs.ElementsAs(ctx, &monitorIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.MonitorIDs = monitorIDs
	}

	if !data.SystemIDs.IsNull() {
		var systemIDs []string
		resp.Diagnostics.Append(data.SystemIDs.ElementsAs(ctx, &systemIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.SystemIDs = systemIDs
	}

	if !data.Metrics.IsNull() {
		createReq.Metrics = data.Metrics.ValueString()
	}

	report, err := r.client.CreateReport(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create report, got error: %s", err))
		return
	}

	r.updateModelFromResponse(&data, report)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ReportResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	report, err := r.client.GetReport(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read report, got error: %s", err))
		return
	}

	r.updateModelFromResponse(&data, report)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Reports cannot be updated - all changes require replacement
	// This method should never be called due to RequiresReplace modifiers
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Reports cannot be updated. All configuration changes require replacement.",
	)
}

func (r *ReportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ReportResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteReport(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete report, got error: %s", err))
		return
	}
}

func (r *ReportResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ReportResource) updateModelFromResponse(data *ReportResourceModel, report *client.Report) {
	data.ID = types.StringValue(report.ID)
	data.Name = types.StringValue(report.Name)
	data.ReportType = types.StringValue(report.ReportType)
	data.Format = types.StringValue(report.Format)
	data.StartTime = types.StringValue(report.StartTime)
	data.EndTime = types.StringValue(report.EndTime)
	data.Status = types.StringValue(report.Status)
	data.CreatedAt = types.StringValue(report.CreatedAt)

	if report.FilePath != "" {
		data.FilePath = types.StringValue(report.FilePath)
	}
	if report.CompletedAt != "" {
		data.CompletedAt = types.StringValue(report.CompletedAt)
	}
	if report.Metrics != "" {
		data.Metrics = types.StringValue(report.Metrics)
	}
}

// setRequiresReplace returns a plan modifier that requires replacement for set attributes.
type setRequiresReplacePlanModifier struct{}

func setRequiresReplace() planmodifier.Set {
	return setRequiresReplacePlanModifier{}
}

func (m setRequiresReplacePlanModifier) Description(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (m setRequiresReplacePlanModifier) MarkdownDescription(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (m setRequiresReplacePlanModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if req.StateValue.IsNull() {
		return
	}

	if req.PlanValue.IsUnknown() {
		return
	}

	if !req.PlanValue.Equal(req.StateValue) {
		resp.RequiresReplace = true
	}
}

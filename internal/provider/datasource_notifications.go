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
var _ datasource.DataSource = &NotificationsDataSource{}

func NewNotificationsDataSource() datasource.DataSource {
	return &NotificationsDataSource{}
}

// NotificationsDataSource defines the data source implementation.
type NotificationsDataSource struct {
	client *client.Client
}

// NotificationsDataSourceModel describes the data source data model.
type NotificationsDataSourceModel struct {
	Page          types.Int64             `tfsdk:"page"`
	PageSize      types.Int64             `tfsdk:"page_size"`
	Total         types.Int64             `tfsdk:"total"`
	TotalPages    types.Int64             `tfsdk:"total_pages"`
	Notifications []NotificationItemModel `tfsdk:"notifications"`
}

// NotificationItemModel describes a single notification history record.
type NotificationItemModel struct {
	ID               types.String `tfsdk:"id"`
	MonitorID        types.String `tfsdk:"monitor_id"`
	AlertID          types.String `tfsdk:"alert_id"`
	IncidentID       types.String `tfsdk:"incident_id"`
	NotificationType types.String `tfsdk:"notification_type"`
	EventType        types.String `tfsdk:"event_type"`
	Destination      types.String `tfsdk:"destination"`
	Subject          types.String `tfsdk:"subject"`
	Message          types.String `tfsdk:"message"`
	Status           types.String `tfsdk:"status"`
	ErrorMessage     types.String `tfsdk:"error_message"`
	ResponseCode     types.Int64  `tfsdk:"response_code"`
	SentAt           types.String `tfsdk:"sent_at"`
	CreatedAt        types.String `tfsdk:"created_at"`
}

func (d *NotificationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notifications"
}

func (d *NotificationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to list notification history.",

		Attributes: map[string]schema.Attribute{
			"page": schema.Int64Attribute{
				MarkdownDescription: "The page number. Default is 1.",
				Optional:            true,
			},
			"page_size": schema.Int64Attribute{
				MarkdownDescription: "The page size. Default is 50, max is 100.",
				Optional:            true,
			},
			"total": schema.Int64Attribute{
				MarkdownDescription: "Total number of notifications.",
				Computed:            true,
			},
			"total_pages": schema.Int64Attribute{
				MarkdownDescription: "Total number of pages.",
				Computed:            true,
			},
			"notifications": schema.ListNestedAttribute{
				MarkdownDescription: "List of notifications.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The notification ID.",
							Computed:            true,
						},
						"monitor_id": schema.StringAttribute{
							MarkdownDescription: "The monitor ID.",
							Computed:            true,
						},
						"alert_id": schema.StringAttribute{
							MarkdownDescription: "The alert ID.",
							Computed:            true,
						},
						"incident_id": schema.StringAttribute{
							MarkdownDescription: "The incident ID.",
							Computed:            true,
						},
						"notification_type": schema.StringAttribute{
							MarkdownDescription: "The type of notification.",
							Computed:            true,
						},
						"event_type": schema.StringAttribute{
							MarkdownDescription: "The event type.",
							Computed:            true,
						},
						"destination": schema.StringAttribute{
							MarkdownDescription: "The notification destination.",
							Computed:            true,
						},
						"subject": schema.StringAttribute{
							MarkdownDescription: "The notification subject.",
							Computed:            true,
						},
						"message": schema.StringAttribute{
							MarkdownDescription: "The notification message.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "The notification status.",
							Computed:            true,
						},
						"error_message": schema.StringAttribute{
							MarkdownDescription: "Any error message.",
							Computed:            true,
						},
						"response_code": schema.Int64Attribute{
							MarkdownDescription: "The response code from the notification endpoint.",
							Computed:            true,
						},
						"sent_at": schema.StringAttribute{
							MarkdownDescription: "When the notification was sent.",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "When the notification was created.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *NotificationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NotificationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NotificationsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	page := 1
	pageSize := 50
	if !data.Page.IsNull() {
		page = int(data.Page.ValueInt64())
	}
	if !data.PageSize.IsNull() {
		pageSize = int(data.PageSize.ValueInt64())
	}

	notificationsResp, err := d.client.ListNotificationHistory(ctx, page, pageSize)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list notifications, got error: %s", err))
		return
	}

	data.Page = types.Int64Value(int64(notificationsResp.Page))
	data.PageSize = types.Int64Value(int64(notificationsResp.PageSize))
	data.Total = types.Int64Value(int64(notificationsResp.Total))
	data.TotalPages = types.Int64Value(int64(notificationsResp.Pages))

	data.Notifications = make([]NotificationItemModel, len(notificationsResp.Notifications))
	for i, notification := range notificationsResp.Notifications {
		data.Notifications[i] = NotificationItemModel{
			ID:               types.StringValue(notification.ID),
			NotificationType: types.StringValue(notification.NotificationType),
			EventType:        types.StringValue(notification.EventType),
			Destination:      types.StringValue(notification.Destination),
			Status:           types.StringValue(notification.Status),
			CreatedAt:        types.StringValue(notification.CreatedAt),
		}
		if notification.MonitorID != "" {
			data.Notifications[i].MonitorID = types.StringValue(notification.MonitorID)
		}
		if notification.AlertID != "" {
			data.Notifications[i].AlertID = types.StringValue(notification.AlertID)
		}
		if notification.IncidentID != "" {
			data.Notifications[i].IncidentID = types.StringValue(notification.IncidentID)
		}
		if notification.Subject != "" {
			data.Notifications[i].Subject = types.StringValue(notification.Subject)
		}
		if notification.Message != "" {
			data.Notifications[i].Message = types.StringValue(notification.Message)
		}
		if notification.ErrorMessage != "" {
			data.Notifications[i].ErrorMessage = types.StringValue(notification.ErrorMessage)
		}
		if notification.ResponseCode != 0 {
			data.Notifications[i].ResponseCode = types.Int64Value(int64(notification.ResponseCode))
		}
		if notification.SentAt != "" {
			data.Notifications[i].SentAt = types.StringValue(notification.SentAt)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

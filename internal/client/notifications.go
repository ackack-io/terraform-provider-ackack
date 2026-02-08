// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
)

// ListNotificationHistory retrieves notification history for the authenticated user.
func (c *Client) ListNotificationHistory(ctx context.Context, page, pageSize int) (*ListNotificationHistoryResponse, error) {
	path := "/api/v1/notifications"
	if page > 0 || pageSize > 0 {
		path = fmt.Sprintf("%s?page=%d&pageSize=%d", path, page, pageSize)
	}
	var resp ListNotificationHistoryResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNotificationHistory retrieves a single notification history record.
func (c *Client) GetNotificationHistory(ctx context.Context, id string) (*NotificationHistory, error) {
	var notification NotificationHistory
	if err := c.get(ctx, fmt.Sprintf("/api/v1/notifications/%s", id), &notification); err != nil {
		return nil, err
	}
	return &notification, nil
}

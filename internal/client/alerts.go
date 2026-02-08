// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
)

// CreateAlert creates a new alert.
func (c *Client) CreateAlert(ctx context.Context, req CreateAlertRequest) (*Alert, error) {
	var alert Alert
	if err := c.post(ctx, "/api/v1/alerts", req, &alert); err != nil {
		return nil, err
	}
	return &alert, nil
}

// GetAlert retrieves an alert by ID.
func (c *Client) GetAlert(ctx context.Context, id string) (*Alert, error) {
	var alert Alert
	if err := c.get(ctx, fmt.Sprintf("/api/v1/alerts/%s", id), &alert); err != nil {
		return nil, err
	}
	return &alert, nil
}

// UpdateAlert updates an existing alert.
func (c *Client) UpdateAlert(ctx context.Context, id string, req UpdateAlertRequest) (*Alert, error) {
	var alert Alert
	if err := c.put(ctx, fmt.Sprintf("/api/v1/alerts/%s", id), req, &alert); err != nil {
		return nil, err
	}
	return &alert, nil
}

// DeleteAlert deletes an alert by ID.
func (c *Client) DeleteAlert(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/api/v1/alerts/%s", id))
}

// ListAlerts retrieves all alerts for the authenticated user.
func (c *Client) ListAlerts(ctx context.Context) ([]Alert, error) {
	var resp ListAlertsResponse
	if err := c.get(ctx, "/api/v1/alerts", &resp); err != nil {
		return nil, err
	}
	return resp.Alerts, nil
}

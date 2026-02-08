// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
)

// CreateReport creates a new report.
func (c *Client) CreateReport(ctx context.Context, req CreateReportRequest) (*Report, error) {
	var report Report
	if err := c.post(ctx, "/api/v1/reports", req, &report); err != nil {
		return nil, err
	}
	return &report, nil
}

// GetReport retrieves a report by ID.
func (c *Client) GetReport(ctx context.Context, id string) (*Report, error) {
	var report Report
	if err := c.get(ctx, fmt.Sprintf("/api/v1/reports/%s", id), &report); err != nil {
		return nil, err
	}
	return &report, nil
}

// DeleteReport deletes a report by ID.
func (c *Client) DeleteReport(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/api/v1/reports/%s", id))
}

// ListReports retrieves all reports for the authenticated user.
func (c *Client) ListReports(ctx context.Context, page, pageSize int) (*ListReportsResponse, error) {
	path := "/api/v1/reports"
	if page > 0 || pageSize > 0 {
		path = fmt.Sprintf("%s?page=%d&pageSize=%d", path, page, pageSize)
	}
	var resp ListReportsResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

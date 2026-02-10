// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
)

// CreateMonitor creates a new monitor.
func (c *Client) CreateMonitor(ctx context.Context, req CreateMonitorRequest) (*Monitor, error) {
	var monitor Monitor
	if err := c.post(ctx, "/api/v1/monitors", req, &monitor); err != nil {
		return nil, err
	}
	return &monitor, nil
}

// GetMonitor retrieves a monitor by ID.
func (c *Client) GetMonitor(ctx context.Context, id string) (*Monitor, error) {
	var monitor Monitor
	if err := c.get(ctx, fmt.Sprintf("/api/v1/monitors/%s", id), &monitor); err != nil {
		return nil, err
	}
	return &monitor, nil
}

// UpdateMonitor updates an existing monitor.
func (c *Client) UpdateMonitor(ctx context.Context, id string, req UpdateMonitorRequest) (*Monitor, error) {
	var monitor Monitor
	if err := c.put(ctx, fmt.Sprintf("/api/v1/monitors/%s", id), req, &monitor); err != nil {
		return nil, err
	}
	return &monitor, nil
}

// DeleteMonitor deletes a monitor by ID.
func (c *Client) DeleteMonitor(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/api/v1/monitors/%s", id))
}

// ListMonitors retrieves all monitors for the authenticated user.
func (c *Client) ListMonitors(ctx context.Context) ([]Monitor, error) {
	var resp ListMonitorsResponse
	if err := c.get(ctx, "/api/v1/monitors", &resp); err != nil {
		return nil, err
	}
	return resp.Monitors, nil
}

// GetMonitorResults retrieves recent check results for a monitor.
func (c *Client) GetMonitorResults(ctx context.Context, id string, limit int) ([]MonitorResult, error) {
	path := fmt.Sprintf("/api/v1/monitors/%s/results", id)
	if limit > 0 {
		path = fmt.Sprintf("%s?limit=%d", path, limit)
	}
	var resp GetResultsResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

// GetMonitorUptime retrieves uptime percentage for a monitor.
func (c *Client) GetMonitorUptime(ctx context.Context, id string, hours int) (*GetUptimeResponse, error) {
	path := fmt.Sprintf("/api/v1/monitors/%s/uptime", id)
	if hours > 0 {
		path = fmt.Sprintf("%s?hours=%d", path, hours)
	}
	var resp GetUptimeResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetMonitorIncidents retrieves recent incidents for a monitor.
func (c *Client) GetMonitorIncidents(ctx context.Context, id string, limit int) ([]Incident, error) {
	path := fmt.Sprintf("/api/v1/monitors/%s/incidents", id)
	if limit > 0 {
		path = fmt.Sprintf("%s?limit=%d", path, limit)
	}
	var resp GetIncidentsResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Incidents, nil
}

// GetMonitorHealth retrieves health information for a specific monitor.
func (c *Client) GetMonitorHealth(ctx context.Context, id string) (*MonitorHealthInfo, error) {
	var resp MonitorHealthInfo
	if err := c.get(ctx, fmt.Sprintf("/api/v1/monitors/%s/health", id), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAllMonitorHealth retrieves health information for all monitors.
func (c *Client) GetAllMonitorHealth(ctx context.Context) (*MonitorHealthResponse, error) {
	var resp MonitorHealthResponse
	if err := c.get(ctx, "/api/v1/monitors/health", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

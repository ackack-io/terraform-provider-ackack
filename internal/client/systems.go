// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
)

// CreateSystem creates a new system.
func (c *Client) CreateSystem(ctx context.Context, req CreateSystemRequest) (*System, error) {
	var system System
	if err := c.post(ctx, "/api/v1/systems", req, &system); err != nil {
		return nil, err
	}
	return &system, nil
}

// GetSystem retrieves a system by ID with statistics.
func (c *Client) GetSystem(ctx context.Context, id string) (*SystemWithStats, error) {
	var system SystemWithStats
	if err := c.get(ctx, fmt.Sprintf("/api/v1/systems/%s", id), &system); err != nil {
		return nil, err
	}
	return &system, nil
}

// UpdateSystem updates an existing system.
func (c *Client) UpdateSystem(ctx context.Context, id string, req UpdateSystemRequest) (*System, error) {
	var system System
	if err := c.put(ctx, fmt.Sprintf("/api/v1/systems/%s", id), req, &system); err != nil {
		return nil, err
	}
	return &system, nil
}

// DeleteSystem deletes a system by ID.
func (c *Client) DeleteSystem(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/api/v1/systems/%s", id))
}

// ListSystems retrieves all systems for the authenticated user.
func (c *Client) ListSystems(ctx context.Context) ([]SystemWithStats, error) {
	var resp ListSystemsResponse
	if err := c.get(ctx, "/api/v1/systems", &resp); err != nil {
		return nil, err
	}
	return resp.Systems, nil
}

// AddMonitorsToSystem adds monitors to a system.
func (c *Client) AddMonitorsToSystem(ctx context.Context, id string, monitorIDs []string) error {
	req := ModifyMonitorsRequest{MonitorIDs: monitorIDs}
	return c.post(ctx, fmt.Sprintf("/api/v1/systems/%s/monitors", id), req, nil)
}

// RemoveMonitorsFromSystem removes monitors from a system.
func (c *Client) RemoveMonitorsFromSystem(ctx context.Context, id string, monitorIDs []string) error {
	req := ModifyMonitorsRequest{MonitorIDs: monitorIDs}
	return c.doRequest(ctx, "DELETE", fmt.Sprintf("/api/v1/systems/%s/monitors", id), req, nil)
}

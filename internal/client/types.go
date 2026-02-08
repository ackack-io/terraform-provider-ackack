// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

// Monitor represents a monitor configuration.
type Monitor struct {
	ID               string  `json:"id,omitempty"`
	UserID           string  `json:"user_id,omitempty"`
	Name             string  `json:"name,omitempty"`
	Type             string  `json:"type,omitempty"`
	IsEnabled        bool    `json:"is_enabled,omitempty"`
	FrequencySeconds int     `json:"frequency_seconds,omitempty"`
	TimeoutMs        int     `json:"timeout_ms,omitempty"`
	Retries          int     `json:"retries,omitempty"`
	GeneralRegion    string  `json:"general_region,omitempty"`
	SpecificRegion   string  `json:"specific_region,omitempty"`
	Status           string  `json:"status,omitempty"`
	UptimePercentage float64 `json:"uptime_percentage,omitempty"`
	LastChecked      string  `json:"last_checked,omitempty"`
	CreatedAt        string  `json:"created_at,omitempty"`
	UpdatedAt        string  `json:"updated_at,omitempty"`

	// HTTP specific
	URL                string `json:"url,omitempty"`
	ExpectedStatusCode int    `json:"expected_status_code,omitempty"`
	ValidateStatus     bool   `json:"validate_status,omitempty"`
	ValidateBody       bool   `json:"validate_body,omitempty"`
	BodyPattern        string `json:"body_pattern,omitempty"`
	Headers            string `json:"headers,omitempty"`

	// DNS specific
	DNSRecordType string `json:"dns_record_type,omitempty"`
	ExpectedValue string `json:"expected_value,omitempty"`
	Nameserver    string `json:"nameserver,omitempty"`

	// TCP specific
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`

	// SSL specific
	Domain                   string `json:"domain,omitempty"`
	CheckExpirationThreshold bool   `json:"check_expiration_threshold,omitempty"`
	ExpirationThreshold      int    `json:"expiration_threshold,omitempty"`
	CheckProtocolVersion     bool   `json:"check_protocol_version,omitempty"`
	MinimumProtocol          string `json:"minimum_protocol,omitempty"`
}

// CreateMonitorRequest is the request body for creating a monitor.
type CreateMonitorRequest struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	IsEnabled        *bool  `json:"is_enabled,omitempty"`
	FrequencySeconds int    `json:"frequency_seconds,omitempty"`
	TimeoutMs        int    `json:"timeout_ms,omitempty"`
	Retries          int    `json:"retries,omitempty"`
	GeneralRegion    string `json:"general_region,omitempty"`
	SpecificRegion   string `json:"specific_region,omitempty"`

	// HTTP specific
	URL                string `json:"url,omitempty"`
	ExpectedStatusCode int    `json:"expected_status_code,omitempty"`
	ValidateStatus     *bool  `json:"validate_status,omitempty"`
	ValidateBody       *bool  `json:"validate_body,omitempty"`
	BodyPattern        string `json:"body_pattern,omitempty"`
	Headers            string `json:"headers,omitempty"`

	// DNS specific
	DNSRecordType string `json:"dns_record_type,omitempty"`
	ExpectedValue string `json:"expected_value,omitempty"`
	Nameserver    string `json:"nameserver,omitempty"`

	// TCP specific
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`

	// SSL specific
	Domain                   string `json:"domain,omitempty"`
	CheckExpirationThreshold *bool  `json:"check_expiration_threshold,omitempty"`
	ExpirationThreshold      int    `json:"expiration_threshold,omitempty"`
	CheckProtocolVersion     *bool  `json:"check_protocol_version,omitempty"`
	MinimumProtocol          string `json:"minimum_protocol,omitempty"`
}

// UpdateMonitorRequest is the request body for updating a monitor.
type UpdateMonitorRequest struct {
	Name             string `json:"name,omitempty"`
	Type             string `json:"type,omitempty"`
	IsEnabled        *bool  `json:"is_enabled,omitempty"`
	FrequencySeconds int    `json:"frequency_seconds,omitempty"`
	TimeoutMs        int    `json:"timeout_ms,omitempty"`
	Retries          int    `json:"retries,omitempty"`
	GeneralRegion    string `json:"general_region,omitempty"`
	SpecificRegion   string `json:"specific_region,omitempty"`

	// HTTP specific
	URL                string `json:"url,omitempty"`
	ExpectedStatusCode int    `json:"expected_status_code,omitempty"`
	ValidateStatus     *bool  `json:"validate_status,omitempty"`
	ValidateBody       *bool  `json:"validate_body,omitempty"`
	BodyPattern        string `json:"body_pattern,omitempty"`
	Headers            string `json:"headers,omitempty"`

	// DNS specific
	DNSRecordType string `json:"dns_record_type,omitempty"`
	ExpectedValue string `json:"expected_value,omitempty"`
	Nameserver    string `json:"nameserver,omitempty"`

	// TCP specific
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`

	// SSL specific
	Domain                   string `json:"domain,omitempty"`
	CheckExpirationThreshold *bool  `json:"check_expiration_threshold,omitempty"`
	ExpirationThreshold      int    `json:"expiration_threshold,omitempty"`
	CheckProtocolVersion     *bool  `json:"check_protocol_version,omitempty"`
	MinimumProtocol          string `json:"minimum_protocol,omitempty"`
}

// ListMonitorsResponse is the response for listing monitors.
type ListMonitorsResponse struct {
	Monitors []Monitor `json:"monitors"`
	Total    int       `json:"total"`
}

// Alert represents an alert configuration.
type Alert struct {
	ID                 string `json:"id,omitempty"`
	UserID             string `json:"user_id,omitempty"`
	MonitorID          string `json:"monitor_id,omitempty"`
	Type               string `json:"type,omitempty"`
	Target             string `json:"target,omitempty"`
	IsEnabled          bool   `json:"is_enabled,omitempty"`
	TriggerThreshold   int    `json:"trigger_threshold,omitempty"`
	RecoveryThreshold  int    `json:"recovery_threshold,omitempty"`
	MinIntervalMinutes int    `json:"min_interval_minutes,omitempty"`
	CustomMessage      string `json:"custom_message,omitempty"`
	IncludeDetails     bool   `json:"include_details,omitempty"`
	LastTriggeredAt    string `json:"last_triggered_at,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
}

// CreateAlertRequest is the request body for creating an alert.
type CreateAlertRequest struct {
	MonitorID          string `json:"monitor_id"`
	Type               string `json:"type"`
	Target             string `json:"target"`
	IsEnabled          *bool  `json:"is_enabled,omitempty"`
	TriggerThreshold   int    `json:"trigger_threshold,omitempty"`
	RecoveryThreshold  int    `json:"recovery_threshold,omitempty"`
	MinIntervalMinutes int    `json:"min_interval_minutes,omitempty"`
	CustomMessage      string `json:"custom_message,omitempty"`
	IncludeDetails     *bool  `json:"include_details,omitempty"`
}

// UpdateAlertRequest is the request body for updating an alert.
type UpdateAlertRequest struct {
	Target             string `json:"target,omitempty"`
	IsEnabled          *bool  `json:"is_enabled,omitempty"`
	TriggerThreshold   int    `json:"trigger_threshold,omitempty"`
	RecoveryThreshold  int    `json:"recovery_threshold,omitempty"`
	MinIntervalMinutes int    `json:"min_interval_minutes,omitempty"`
	CustomMessage      string `json:"custom_message,omitempty"`
	IncludeDetails     *bool  `json:"include_details,omitempty"`
}

// ListAlertsResponse is the response for listing alerts.
type ListAlertsResponse struct {
	Alerts []Alert `json:"alerts"`
}

// ExternalLink represents an external link on a system.
type ExternalLink struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// System represents a system grouping of monitors.
type System struct {
	ID            string         `json:"id,omitempty"`
	UserID        string         `json:"user_id,omitempty"`
	Name          string         `json:"name,omitempty"`
	Description   string         `json:"description,omitempty"`
	Priority      string         `json:"priority,omitempty"`
	Status        string         `json:"status,omitempty"`
	ExternalLinks []ExternalLink `json:"external_links,omitempty"`
	CreatedAt     string         `json:"created_at,omitempty"`
	UpdatedAt     string         `json:"updated_at,omitempty"`
}

// SystemWithStats represents a system with aggregated statistics.
type SystemWithStats struct {
	ID            string         `json:"id,omitempty"`
	UserID        string         `json:"user_id,omitempty"`
	Name          string         `json:"name,omitempty"`
	Description   string         `json:"description,omitempty"`
	Priority      string         `json:"priority,omitempty"`
	Status        string         `json:"status,omitempty"`
	ExternalLinks []ExternalLink `json:"external_links,omitempty"`
	MonitorCount  int            `json:"monitor_count,omitempty"`
	HealthyCount  int            `json:"healthy_count,omitempty"`
	DegradedCount int            `json:"degraded_count,omitempty"`
	ErrorCount    int            `json:"error_count,omitempty"`
	WarningCount  int            `json:"warning_count,omitempty"`
	OverallUptime float64        `json:"overall_uptime,omitempty"`
	CreatedAt     string         `json:"created_at,omitempty"`
	UpdatedAt     string         `json:"updated_at,omitempty"`
}

// CreateSystemRequest is the request body for creating a system.
type CreateSystemRequest struct {
	Name          string         `json:"name"`
	Description   string         `json:"description,omitempty"`
	Priority      string         `json:"priority,omitempty"`
	ExternalLinks []ExternalLink `json:"external_links,omitempty"`
	MonitorIDs    []string       `json:"monitor_ids"`
}

// UpdateSystemRequest is the request body for updating a system.
type UpdateSystemRequest struct {
	Name          string         `json:"name,omitempty"`
	Description   string         `json:"description,omitempty"`
	Priority      string         `json:"priority,omitempty"`
	ExternalLinks []ExternalLink `json:"external_links,omitempty"`
}

// ListSystemsResponse is the response for listing systems.
type ListSystemsResponse struct {
	Systems []SystemWithStats `json:"systems"`
	Total   int               `json:"total"`
}

// ModifyMonitorsRequest is the request for adding/removing monitors from a system.
type ModifyMonitorsRequest struct {
	MonitorIDs []string `json:"monitor_ids"`
}

// Report represents a generated report.
type Report struct {
	ID            string   `json:"id,omitempty"`
	UserID        string   `json:"user_id,omitempty"`
	Name          string   `json:"name,omitempty"`
	ReportType    string   `json:"report_type,omitempty"`
	Format        string   `json:"format,omitempty"`
	Status        string   `json:"status,omitempty"`
	StartTime     string   `json:"start_time,omitempty"`
	EndTime       string   `json:"end_time,omitempty"`
	MonitorIDs    []string `json:"monitor_ids,omitempty"`
	Metrics       string   `json:"metrics,omitempty"`
	Data          string   `json:"data,omitempty"`
	FilePath      string   `json:"file_path,omitempty"`
	FileSizeBytes int      `json:"file_size_bytes,omitempty"`
	ErrorMessage  string   `json:"error_message,omitempty"`
	CompletedAt   string   `json:"completed_at,omitempty"`
	CreatedAt     string   `json:"created_at,omitempty"`
}

// CreateReportRequest is the request body for creating a report.
type CreateReportRequest struct {
	Name       string   `json:"name"`
	ReportType string   `json:"report_type"`
	Format     string   `json:"format"`
	StartTime  string   `json:"start_time"`
	EndTime    string   `json:"end_time"`
	MonitorIDs []string `json:"monitor_ids,omitempty"`
	SystemIDs  []string `json:"system_ids,omitempty"`
	Metrics    string   `json:"metrics,omitempty"`
}

// ListReportsResponse is the response for listing reports.
type ListReportsResponse struct {
	Reports  []Report `json:"reports"`
	Total    int      `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
	Pages    int      `json:"pages"`
}

// MonitorResult represents a single check result.
type MonitorResult struct {
	ID                       int    `json:"id,omitempty"`
	MonitorID                string `json:"monitor_id,omitempty"`
	Status                   string `json:"status,omitempty"`
	ResponseTime             int    `json:"response_time,omitempty"`
	ResponseSizeBytes        int    `json:"response_size_bytes,omitempty"`
	Timestamp                string `json:"timestamp,omitempty"`
	Region                   string `json:"region,omitempty"`
	WorkerID                 string `json:"worker_id,omitempty"`
	Message                  string `json:"message,omitempty"`
	ErrorType                string `json:"error_type,omitempty"`
	StatusCode               int    `json:"status_code,omitempty"`
	DNSResponse              string `json:"dns_response,omitempty"`
	TLSVersion               string `json:"tls_version,omitempty"`
	CertificateExpirationDays int   `json:"certificate_expiration_days,omitempty"`
}

// GetResultsResponse is the response for getting monitor results.
type GetResultsResponse struct {
	Results []MonitorResult `json:"results"`
	Total   int             `json:"total"`
}

// GetUptimeResponse is the response for getting monitor uptime.
type GetUptimeResponse struct {
	MonitorID string  `json:"monitor_id"`
	Hours     int     `json:"hours"`
	Uptime    float64 `json:"uptime"`
}

// Incident represents a monitor incident.
type Incident struct {
	ID              string `json:"id,omitempty"`
	MonitorID       string `json:"monitor_id,omitempty"`
	Status          string `json:"status,omitempty"`
	Severity        string `json:"severity,omitempty"`
	Summary         string `json:"summary,omitempty"`
	Details         string `json:"details,omitempty"`
	FirstErrorID    string `json:"first_error_id,omitempty"`
	StartedAt       string `json:"started_at,omitempty"`
	ResolvedAt      string `json:"resolved_at,omitempty"`
	DurationSeconds int    `json:"duration_seconds,omitempty"`
	Notified        bool   `json:"notified,omitempty"`
}

// GetIncidentsResponse is the response for getting monitor incidents.
type GetIncidentsResponse struct {
	Incidents []Incident `json:"incidents"`
}

// MonitorHealthInfo represents health information for a single monitor.
type MonitorHealthInfo struct {
	MonitorID       string  `json:"monitor_id,omitempty"`
	MonitorName     string  `json:"monitor_name,omitempty"`
	IsInFlight      bool    `json:"is_in_flight,omitempty"`
	InFlightSeconds float64 `json:"in_flight_seconds,omitempty"`
	Throttled       bool    `json:"throttled,omitempty"`
	ThrottleReason  string  `json:"throttle_reason,omitempty"`
	DampeningLevel  int     `json:"dampening_level,omitempty"`
	DampeningName   string  `json:"dampening_name,omitempty"`
	DampeningReason string  `json:"dampening_reason,omitempty"`
	FailureRate     float64 `json:"failure_rate,omitempty"`
	P95LatencyMs    int     `json:"p95_latency_ms,omitempty"`
	StuckCount      int     `json:"stuck_count,omitempty"`
}

// UserHealthSummary represents health summary for a user.
type UserHealthSummary struct {
	Plan           string `json:"plan,omitempty"`
	InFlightCount  int    `json:"in_flight_count,omitempty"`
	InFlightLimit  int    `json:"in_flight_limit,omitempty"`
	AtLimit        bool   `json:"at_limit,omitempty"`
	ThrottledCount int    `json:"throttled_count,omitempty"`
	DampenedCount  int    `json:"dampened_count,omitempty"`
}

// MonitorHealthResponse is the response for getting all monitor health.
type MonitorHealthResponse struct {
	Monitors []MonitorHealthInfo `json:"monitors"`
	User     UserHealthSummary   `json:"user"`
}

// NotificationHistory represents a notification history record.
type NotificationHistory struct {
	ID               string `json:"id,omitempty"`
	UserID           string `json:"user_id,omitempty"`
	MonitorID        string `json:"monitor_id,omitempty"`
	AlertID          string `json:"alert_id,omitempty"`
	IncidentID       string `json:"incident_id,omitempty"`
	NotificationType string `json:"notification_type,omitempty"`
	EventType        string `json:"event_type,omitempty"`
	Destination      string `json:"destination,omitempty"`
	Subject          string `json:"subject,omitempty"`
	Message          string `json:"message,omitempty"`
	Details          string `json:"details,omitempty"`
	Status           string `json:"status,omitempty"`
	ErrorMessage     string `json:"error_message,omitempty"`
	ResponseCode     int    `json:"response_code,omitempty"`
	ResponseBody     string `json:"response_body,omitempty"`
	DeliveryAttempts int    `json:"delivery_attempts,omitempty"`
	SentAt           string `json:"sent_at,omitempty"`
	LastAttemptAt    string `json:"last_attempt_at,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
}

// ListNotificationHistoryResponse is the response for listing notification history.
type ListNotificationHistoryResponse struct {
	Notifications []NotificationHistory `json:"notifications"`
	Total         int                   `json:"total"`
	Page          int                   `json:"page"`
	PageSize      int                   `json:"pageSize"`
	Pages         int                   `json:"pages"`
}

// ErrorResponse is the API error response.
type ErrorResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

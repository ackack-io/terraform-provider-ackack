resource "ackack_report" "monthly_uptime" {
  name        = "Monthly Uptime Report"
  report_type = "uptime"
  format      = "pdf"
  start_time  = "2024-01-01T00:00:00Z"
  end_time    = "2024-01-31T23:59:59Z"

  monitor_ids = [
    ackack_monitor.website.id,
    ackack_monitor.dns.id,
  ]
}

resource "ackack_report" "incidents" {
  name        = "Incident Report"
  report_type = "incidents"
  format      = "json"
  start_time  = "2024-01-01T00:00:00Z"
  end_time    = "2024-01-31T23:59:59Z"
}

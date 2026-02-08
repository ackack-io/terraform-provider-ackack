# Email Alert
resource "ackack_alert" "email" {
  monitor_id          = ackack_monitor.website.id
  type                = "email"
  target              = "alerts@example.com"
  is_enabled          = true
  trigger_threshold   = 2
  recovery_threshold  = 2
  min_interval_minutes = 15
  include_details     = true
}

# Slack Alert
resource "ackack_alert" "slack" {
  monitor_id = ackack_monitor.website.id
  type       = "slack"
  target     = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
  is_enabled = true
}

# Webhook Alert
resource "ackack_alert" "webhook" {
  monitor_id     = ackack_monitor.website.id
  type           = "webhook"
  target         = "https://api.example.com/webhooks/monitoring"
  is_enabled     = true
  custom_message = "Monitor {{monitor_name}} is {{status}}"
}

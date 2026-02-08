# HTTP Monitor
resource "ackack_monitor" "website" {
  name              = "Website Monitor"
  type              = "http"
  url               = "https://example.com"
  frequency_seconds = 60
  timeout_ms        = 10000
  is_enabled        = true

  validate_status      = true
  expected_status_code = 200
}

# DNS Monitor
resource "ackack_monitor" "dns" {
  name              = "DNS Monitor"
  type              = "dns"
  url               = "example.com"
  dns_record_type   = "A"
  expected_value    = "93.184.216.34"
  frequency_seconds = 300
}

# SSL Monitor
resource "ackack_monitor" "ssl" {
  name                       = "SSL Certificate Monitor"
  type                       = "ssl"
  domain                     = "example.com"
  check_expiration_threshold = true
  expiration_threshold       = 30
  frequency_seconds          = 3600
}

# TCP Monitor
resource "ackack_monitor" "tcp" {
  name              = "TCP Port Monitor"
  type              = "tcp"
  host              = "example.com"
  port              = 443
  frequency_seconds = 60
  timeout_ms        = 5000
}

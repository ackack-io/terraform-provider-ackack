resource "ackack_system" "production" {
  name        = "Production System"
  description = "All production services"
  priority    = "critical"

  monitor_ids = [
    ackack_monitor.website.id,
    ackack_monitor.dns.id,
    ackack_monitor.ssl.id,
  ]

  external_links {
    name = "Dashboard"
    url  = "https://dashboard.example.com"
  }

  external_links {
    name = "Documentation"
    url  = "https://docs.example.com"
  }
}

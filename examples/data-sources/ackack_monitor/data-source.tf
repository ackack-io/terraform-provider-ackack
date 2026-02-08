data "ackack_monitor" "example" {
  id = "mon_abc123"
}

output "monitor_status" {
  value = data.ackack_monitor.example.status
}

output "uptime_percentage" {
  value = data.ackack_monitor.example.uptime_percentage
}

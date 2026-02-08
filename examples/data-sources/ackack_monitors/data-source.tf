data "ackack_monitors" "all" {}

output "total_monitors" {
  value = length(data.ackack_monitors.all.monitors)
}

output "enabled_monitors" {
  value = [
    for m in data.ackack_monitors.all.monitors : m.name
    if m.is_enabled
  ]
}

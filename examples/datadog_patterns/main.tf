resource "datadog_monitor" "valid_monitor" {
  name    = "Valid Monitor"
  type    = "metric alert"
  message = "Valid monitor with correct patterns"
  query   = "avg(last_5m):avg:system.cpu.user{*} > 0.8"

  tags = [
    "Environment:prod",
    "Team:platform-team",
    "Service:web-api",
    "Version:1.2.3"
  ]
}

resource "datadog_dashboard" "invalid_dashboard" {
  title       = "Invalid Dashboard"
  description = "Dashboard with pattern violations"
  layout_type = "ordered"

  tags = [
    "Environment:PRODUCTION",  # Should fail - uppercase
    "Team:Platform Team",      # Should fail - contains space
    "Service:WebAPI",          # Should fail - uppercase
    "Version:latest"           # Should fail - not semantic version
  ]

  widget {
    timeseries_definition {
      request {
        query {
          metric_query {
            data_source = "metrics"
            query       = "avg:system.cpu.user{*}"
            name        = "cpu_query"
          }
        }
      }
    }
  }
}

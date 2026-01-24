resource "datadog_monitor" "cpu_monitor" {
  name    = "High CPU Usage"
  type    = "metric alert"
  message = "CPU usage is high"
  query   = "avg(last_5m):avg:system.cpu.user{*} > 0.8"

  tags = [
    "Name:cpu-monitor",
    "Environment:production",
    "Owner:platform-team",
    "Project:monitoring"
  ]
}

resource "datadog_dashboard" "main_dashboard" {
  title       = "Main Dashboard"
  description = "Main monitoring dashboard"
  layout_type = "ordered"

  tags = [
    "Name:main-dashboard",
    "Environment:production",
    "Owner:platform-team",
    "Project:monitoring"
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

resource "datadog_synthetics_test" "api_test" {
  name    = "API Health Check"
  type    = "api"
  subtype = "http"
  status  = "live"
  locations = ""

  request_definition {
    method = "GET"
    url    = "https://api.example.com/health"
  }

  tags = [
    "Name:api-health-check",
    "Environment:production",
    "Owner:platform-team",
    "Project:monitoring"
  ]
}

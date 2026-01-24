resource "datadog_monitor" "cpu_monitor" {
  name    = "High CPU Usage"
  type    = "metric alert"
  message = "CPU usage is high"
  query   = "avg(last_5m):avg:system.cpu.user{*} > 0.8"

  tags = [
    "Environment:production",
    "Team:platform",
    "Service:web-api"
  ]
}

resource "datadog_dashboard" "main_dashboard" {
  title       = "Main Dashboard"
  description = "Main monitoring dashboard"
  layout_type = "ordered"

  tags = [
    "Environment:production",
    "Team:platform",
    "Service:web-api"
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

  request_definition {
    method = "GET"
    url    = "https://api.example.com/health"
  }

  tags = [
    "Environment:production",
    "Team:platform",
    "Service:web-api"
  ]
}

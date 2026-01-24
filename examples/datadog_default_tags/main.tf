provider "datadog" {
  api_key = var.datadog_api_key
  app_key = var.datadog_app_key

  default_tags {
    tags = {
      Environment = "production"
      Team        = "platform"
    }
  }
}

resource "datadog_monitor" "cpu_monitor" {
  name    = "High CPU Usage"
  type    = "metric alert"
  message = "CPU usage is high"
  query   = "avg(last_5m):avg:system.cpu.user{*} > 0.8"

  tags = [
    "Service:web-api",
    "Priority:high"
  ]
}

resource "datadog_service_level_objective" "api_slo" {
  name        = "API Availability SLO"
  type        = "monitor"
  description = "API should be available 99.9% of the time"

  monitor_ids = [datadog_monitor.cpu_monitor.id]

  thresholds {
    timeframe = "7d"
    target    = 99.9
    warning   = 99.95
  }

  tags = [
    "Service:web-api",
    "Type:availability"
  ]
}

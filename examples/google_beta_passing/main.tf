provider "google-beta" {
  project = "my-project-id"
  region  = "us-central1"

  default_labels = {
    Environment = "beta"
    Owner       = "team-beta"
  }
}

# Example with all required labels using google-beta provider
resource "google_compute_instance" "beta_compliant" {
  provider     = google-beta
  name         = "beta-instance"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }

  labels = {
    Name    = "beta-instance"
    Project = "beta-demo"
  }
}

# Example with beta-specific features
resource "google_storage_bucket" "beta_compliant" {
  provider = google-beta
  name     = "my-beta-bucket"
  location = "US"

  labels = {
    Name    = "my-beta-bucket"
    Project = "beta-demo"
  }
}

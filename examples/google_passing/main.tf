provider "google" {
  project = "my-project-id"
  region  = "us-central1"

  default_labels = {
    Environment = "dev"
    Owner       = "team-a"
  }
}

# Example with all required labels
resource "google_compute_instance" "compliant" {
  name         = "compliant-instance"
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
    Name    = "compliant-instance"
    Project = "demo"
  }
}

# Example with all required labels from provider defaults and resource
resource "google_storage_bucket" "compliant" {
  name     = "my-bucket"
  location = "US"

  labels = {
    Name    = "my-bucket"
    Project = "demo"
  }
}

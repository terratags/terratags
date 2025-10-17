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

# Example with missing labels
resource "google_storage_bucket" "non_compliant" {
  name     = "my-bucket"
  location = "US"

  labels = {
    Name = "my-bucket"
  }
}

# Example with provider default_labels only
resource "google_compute_disk" "with_defaults" {
  name = "test-disk"
  type = "pd-ssd"
  zone = "us-central1-a"
  size = 10

  labels = {
    Name    = "test-disk"
    Project = "demo"
  }
}

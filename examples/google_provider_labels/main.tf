terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }
}

provider "google" {
  project = "my-project-id"
  region  = "us-central1"
  
  default_labels= {
    labels = {
      Environment = "dev"
      Owner       = "team-a"
      Project     = "demo"
    }
  }
}

resource "google_compute_instance" "example" {
  name         = "example-instance"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-10"
    }
  }

  network_interface {
    network = "default"
  }

  # Only need to specify Name label, as other required labels come from default_labels
  labels = {
    Name = "example-instance"
  }
}

resource "google_storage_bucket" "example" {
  name          = "example-bucket"
  location      = "US"
  force_destroy = true

  # Adding Name label, other required labels come from default_labels
  labels = {
    Name = "example-bucket"
  }
}

# Resource without labels - will be flagged as non-compliant
resource "google_compute_firewall" "example" {
  name    = "example-firewall"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["80", "443"]
  }

  source_ranges = ["0.0.0.0/0"]
}
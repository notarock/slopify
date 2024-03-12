terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "5.7.0"
    }
  }
}

provider "google" {
  project     = "slopify-414804"
  region      = "northamerica-northeast1"
}



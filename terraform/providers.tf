terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "5.7.0"
    }
  }
}

provider "google" {
  project     = "sludger"
  region      = "northamerica-northeast1"
}



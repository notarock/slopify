resource "google_storage_bucket" "public_bucket" {
  name     = "sludger-temp"
  location = "northamerica-northeast1"

  force_destroy = true

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }

  cors {
    origin          = ["*"]
    method          = ["GET"]
    response_header = ["Content-Type"]
    max_age_seconds = 3600
  }

  lifecycle_rule {
    condition {
      age = "3"
    }
    action {
      type = "Delete"
    }
  }
}

resource "google_storage_bucket_iam_binding" "public_access" {
  bucket = google_storage_bucket.public_bucket.name
  role   = "roles/storage.objectViewer"

  members = [
    "allUsers",
  ]
}

output "bucket_url" {
  value = google_storage_bucket.public_bucket.url
}

output "bucket_name" {
  value = google_storage_bucket.public_bucket.name
}

resource "google_project_service" "this" {
  for_each = local.stack
  service = "cloudfunctions.googleapis.com"
  project = data.google_project.this[each.key].project_id
}

resource "google_project_service" "cloud_run_api" {
  project = var.project_id
  service = "run.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "cloudbuild" {
  project = var.project_id
  service = "cloudbuild.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "artifactregistry" {
  project = var.project_id
  service = "artifactregistry.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "firestore_api" {
  project = var.project_id
  service = "firestore.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "pubsub_api" {
  project = var.project_id
  service = "pubsub.googleapis.com"
  disable_dependent_services = true
}

resource "google_project_service" "scheduler_api" {
  project = var.project_id
  service = "cloudscheduler.googleapis.com"
  disable_dependent_services = true
}

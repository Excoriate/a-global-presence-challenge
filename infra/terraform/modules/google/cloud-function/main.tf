resource "random_id" "this" {
  for_each = local.stack
  byte_length = 6
}

resource "random_id" "sa" {
  for_each = local.stack
  byte_length = 3
}

resource "google_service_account" "this" {
  for_each = local.stack
  account_id   = substr(lower(format("sa-%s-%s", var.function_name, random_id.this[each.key].hex)), 0, 29)
  display_name = format("Service Account for %s", local.resource_name)
  project = data.google_project.this[each.key].project_id
}

resource "google_project_iam_binding" "firestore_editor" {
  for_each = local.stack
  role     = "roles/datastore.user"
  members  = [
    "serviceAccount:${google_service_account.this[each.key].email}",
  ]
  project  = data.google_project.this[each.key].project_id
}

resource "google_service_account_iam_member" "service_account_user" {
  for_each = local.stack
  service_account_id = google_service_account.this[each.key].id
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.this[each.key].email}"
}

resource "google_project_iam_binding" "binding_service_account" {
  for_each = local.stack
  role     = "roles/cloudfunctions.developer"
  members  = [
    "serviceAccount:${google_service_account.this[each.key].email}",
  ]
  project  = data.google_project.this[each.key].project_id
}

resource "google_project_iam_binding" "this" {
  for_each = local.stack
  role     = "roles/iam.serviceAccountUser"
  members  = [
    "serviceAccount:${google_service_account.this[each.key].email}",
  ]
  project  = data.google_project.this[each.key].project_id
}


resource "time_sleep" "this" {
  for_each = local.stack
  create_duration = "60s"
}

resource "google_cloud_run_service_iam_policy" "public" {
  for_each = local.stack
  location = var.location
  project  = data.google_project.this[each.key].project_id
  service  = google_cloudfunctions2_function.this[each.key].name

  policy_data = data.google_iam_policy.public[each.key].policy_data
}

resource "google_cloudfunctions2_function" "this" {
  for_each = local.stack
  name                  = local.resource_name
  location = var.location
  description = format("Function for %s", local.resource_name)
  build_config {
    runtime               = var.runtime
    entry_point = var.function_name
    source {
      storage_source {
        bucket = google_storage_bucket_object.this[each.key].bucket
        object = google_storage_bucket_object.this[each.key].name
      }
    }

    environment_variables = local.env_vars
  }


  service_config {
    max_instance_count = 1
    available_memory = var.memory
    timeout_seconds = 60
    ingress_settings = "ALLOW_ALL"
    all_traffic_on_latest_revision = true
    service_account_email = google_service_account.this[each.key].email
  }

  project = data.google_project.this[each.key].project_id


  labels = local.labels

  depends_on = [time_sleep.this]

  lifecycle {
    create_before_destroy = false
  }
}

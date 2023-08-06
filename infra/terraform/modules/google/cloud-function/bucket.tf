resource "google_storage_bucket" "this" {
  for_each = local.stack
  name = format("%s-%s", local.resource_name, random_id.this[each.key].hex)
  project = data.google_project.this[each.key].project_id
  location = var.location
  storage_class = "STANDARD"
  force_destroy = true
}

resource "google_storage_bucket_object" "this" {
  for_each = local.stack
  name = format("deploy/%s", format("%s-%s.zip", local.resource_name, data.archive_file.this[each.key].output_md5))
  bucket = google_storage_bucket.this[each.key].name
  source = data.archive_file.this[each.key].output_path
}

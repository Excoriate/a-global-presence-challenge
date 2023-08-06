data "google_project" "this" {
  for_each = local.stack
  project_id = local.project_id
}

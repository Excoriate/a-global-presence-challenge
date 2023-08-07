data "google_project" "this" {
  for_each = local.stack_to_deploy
  project_id = local.project_id
}

data "archive_file" "this" {
  for_each = local.stack_to_deploy
  type        = "zip"
  source_dir= var.source_code_path
  output_path = format("%s/output/%s.zip", path.module, var.function_name)
}

data "google_iam_policy" "public" {
  for_each = local.stack_to_deploy
  binding {
    role    = "roles/run.invoker"
    members = ["allUsers"]
  }
}

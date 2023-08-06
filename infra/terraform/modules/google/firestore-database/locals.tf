locals {
  /*
    * Feature flags
  */
  is_enabled = var.is_enabled
  stack = !local.is_enabled ? {} : {
    id = var.owner
  }

  default_name = format("%s-%s", var.owner, var.environment)
  resource_name = !local.is_enabled ? null : var.resource_friendly_identifier == null ? local.default_name : trimspace(var.resource_friendly_identifier)
  project_id = !local.is_enabled ? null : trimspace(var.project_id)
}

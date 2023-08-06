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

  labels = !local.is_enabled ? {} : {
      owner           = var.owner,
      environment     = var.environment,
      resource_name   = local.resource_name
      project_id    = local.project_id
      resource_type = "stack"
  }

  default_environment_variables = !local.is_enabled ? null : {
    ENVIRONMENT = var.environment
    LOCATION    = var.location
    REGION      = var.region
    IS_ENABLED  = var.is_enabled // Set by default.
  }

  env_vars = merge(local.default_environment_variables, var.environment_variables)
}

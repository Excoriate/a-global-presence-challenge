locals {
  /*
    * Feature flags
  */
  is_enabled = var.is_enabled
  stack = !local.is_enabled ? {} : {
    id = var.owner
  }

  combinations = flatten([
    for s in local.stack : [
      for l in var.location : {
        stack_key = s
        location  = l
      }
    ]
  ])

  stack_to_deploy = {
    for c in local.combinations: "${c["stack_key"]}-${c["location"]}" => c
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

  default_environment_variables = !local.is_enabled ? null : [
    {
      name = "OWNER"
      value = var.owner
    },
    {
      name = "ENVIRONMENT"
      value = var.environment
    },
    {
      name = "RESOURCE_NAME"
      value = local.resource_name
    },
    {
      name = "PROJECT_ID"
      value = local.project_id
    },
    {
      name = "REGION"
      value = var.region
    },
    {
      name = "IS_ENABLED"
      value = tostring(var.is_enabled)
    }
  ]

  access_token_env_var = !local.is_enabled ? null : [{
    name = "ACCESS_TOKEN"
    value = var.access_token_env_var  == null ? "" : var.access_token_env_var
  }]

  env_vars = concat(local.default_environment_variables, var.environment_variables, local.access_token_env_var)

    env_vars_map = !local.is_enabled ? null : {
        for e in local.env_vars : e["name"] => e["value"]
    }
}

locals {
  environment = get_env("TF_VAR_environment")

  terraform_version  = get_env("TERRAGRUNT_CFG_BINARIES_TERRAFORM_VERSION", "1.5.1")
  terragrunt_version = get_env("TERRAGRUNT_CFG_BINARIES_TERRAGRUNT_VERSION", "0.42.8")

  resolved_type  = trimspace(get_env("TERRAGRUNT_CFG_STATE_RUNTIME_TYPE", "infra"))
  resolved_layer = local.resolved_type == "web" ? "frontend" : "backend"
  project_id    = get_env("TF_VAR_project_id")

  location = trimspace(get_env("LOCATION", "europe-west4"))
  region     = trimspace(get_env("REGION", "europe-west4"))

  #######################################################
  # Remote backend configuration
  #######################################################
  key = {
    prefix      = trimspace(get_env("TERRAGRUNT_CFG_STATE_RUNTIME_PREFIX", "OSS"))
    owner       = trimspace(get_env("TERRAGRUNT_CFG_STATE_RUNTIME_PREFIX", "alex"))
    project_id  = local.project_id
    type        = "lab"
    layer       = local.project_id
    domain      = trimspace(get_env("TERRAGRUNT_CFG_STATE_RUNTIME_DOMAIN", "hackattic"))
    environment = local.environment
    region = format("region-%s", local.region)
    location = format("location-%s", local.location)
    base_name   = "terraform.tfstate"
  }

  terraform_state_file_bucket_region       = get_env("TF_STATE_BUCKET_REGION", "eu-central-1")
  terraform_state_file_bucket              = get_env("TF_STATE_BUCKET", format("platform-tfstate-account-%s", local.environment))
  terraform_state_file_lock_dynamodb_table = get_env("TF_STATE_LOCK_TABLE", format("platform-tfstate-account-%s", local.environment))

  key_str = join("/", [
    local.key.prefix,
    local.key.owner,
    local.key.project_id,
    format("domain-%s", local.key.domain),
    format("layer-%s", local.resolved_layer),
    format("type-%s", local.resolved_type),
    local.key.region,
    local.key.location,
    local.deployment_path,
    local.key.environment,
    local.key.base_name
  ])

  deployment_path = "stack-${path_relative_to_include()}"


  #######################################################
  # Expose metadata for troubleshooting purposes.
  #######################################################
  expose_line_separator     = run_cmd("sh", "-c", "echo '================================================================================'")
  expose_terraform_key_path = run_cmd("sh", "-c", format("export TERRAFORM_STATE_KEYPATH=%s; echo Terraform state keypath : $TERRAFORM_STATE_KEYPATH", local.key_str))
  expose_tf_bucket_region   = run_cmd("sh", "-c", format("export TF_STATE_BUCKET_REGION=%s; echo Terraform state bucket region : $TF_STATE_BUCKET_REGION", local.terraform_state_file_bucket_region))
  expose_tf_state_bucket    = run_cmd("sh", "-c", format("export TF_STATE_BUCKET=%s; echo Terraform state bucket : $TF_STATE_BUCKET", local.terraform_state_file_bucket))
  expose_environment        = run_cmd("sh", "-c", format("export ENVIRONMENT=%s; echo Environment : $ENVIRONMENT", local.environment))
  expose_domain         = run_cmd("sh", "-c", format("export TSN_DOMAIN=%s; echo TSN Domain : $TSN_DOMAIN", local.key.domain))
  expose_layer          = run_cmd("sh", "-c", format("export TSN_LAYER=%s; echo TSN Layer : $TSN_LAYER", local.key.layer))
  expose_project_id     = run_cmd("sh", "-c", format("export PROJECT_ID=%s; echo Project ID : $PROJECT_ID", local.project_id))
  expose_region         = run_cmd("sh", "-c", format("export REGION=%s; echo Region : $REGION", local.region))
    expose_location       = run_cmd("sh", "-c", format("export LOCATION=%s; echo Location : $LOCATION", local.location))
}

terraform {
  extra_arguments "optional_vars" {
    commands = [
      "apply",
      "destroy",
      "plan",
    ]

    required_var_files = [
      "${get_terragrunt_dir()}/envs/common.tfvars",
      "${get_terragrunt_dir()}/envs/${local.environment}.tfvars",
    ]

    optional_var_files = [
      "${get_terragrunt_dir()}/envs/${local.environment}.tfvars",
      "${get_terragrunt_dir()}/envs/${local.environment}-${local.region}.tfvars",
      "${get_terragrunt_dir()}/envs/common-${local.region}.tfvars",
    ]
  }

  extra_arguments "disable_input" {
    commands  = get_terraform_commands_that_need_input()
    arguments = ["-input=false"]
  }

  after_hook "clean_cache_after_apply" {
    commands = ["apply"]
    execute  = ["rm", "-rf", ".terragrunt-cache"]
  }

  after_hook "remove_auto_generated_backend" {
    commands = ["apply"]
    execute  = ["rm", "-rf", "backend.tf"]
  }

  after_hook "remove_auto_generated_provider" {
    commands = ["apply"]
    execute  = ["rm", "-rf", "provider.tf"]
  }
}


generate "terraform_version" {
  path              = ".terraform-version"
  if_exists         = "overwrite"
  disable_signature = true

  contents = <<-EOF
    ${local.terraform_version}
  EOF
}

generate "terragrunt_version" {
  path              = ".terragrunt-version"
  if_exists         = "overwrite"
  disable_signature = true

  contents = <<-EOF
    ${local.terragrunt_version}
  EOF
}


generate "providers" {
  path      = "providers.tf"
  if_exists = "overwrite_terragrunt"

  contents = templatefile("${get_repo_root()}/infra/terraform/stacks/providers.tf.tmpl", {
  })
}

#######################################################
# Common inputs passed to all modules
#######################################################
inputs = {
  environment = local.environment
  owner      = local.key.owner
  project_id  = local.project_id
}

#######################################################
# Remote backend configuration
#######################################################
remote_state {
  backend = "s3"

  generate = {
    path      = "backend.tf"
    if_exists = "overwrite"
  }

  config = {
    disable_bucket_update = true
    encrypt               = true

    region         = local.terraform_state_file_bucket_region
    dynamodb_table = local.terraform_state_file_lock_dynamodb_table
    bucket         = local.terraform_state_file_bucket

    key = local.key_str
  }
}

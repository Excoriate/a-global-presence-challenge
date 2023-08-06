include "root" {
  path           = find_in_parent_folders()
  merge_strategy = "deep"
}

include "parent" {
  path           = "${get_terragrunt_dir()}/../../_env/_cfg/_metadata.hcl"
  expose         = true
  merge_strategy = "deep"
}


locals {
  global_tags = include.parent.locals.tags
  local_tags = {
  }
}

terraform {
  source = "../../../modules//google/cloud-function"
}

inputs = {
  tags = merge(local.global_tags, local.local_tags)
  source_code_path= format("${get_repo_root()}/functions/trigger")
}

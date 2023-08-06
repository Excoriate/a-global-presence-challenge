locals {
  tags = {
    ManagedBy         = "Terraform"
    OrchestratedBy    = "Terragrunt"
    Author            = get_env("TERRAGRUNT_CFG_METADATA_RUNTIME_AUTHOR", "alex@ideaup.cl")
    Owner             = "github.com/Excoriate"
    Type              = "infrastructure"
    Maintainer        = "AlexTorres"
    Project           = get_env("TERRAGRUNT_CFG_METADATA_RUNTIME_PROJECT", "OSS")
    Environment       = get_env("TF_VAR_environment")
    Owner             = trimspace(get_env("TERRAGRUNT_CFG_STATE_RUNTIME_PREFIX", "AlexTorres")) //tsn/other
    Type              = "infrastructure"
    Domain            = trimspace(get_env("TERRAGRUNT_CFG_STATE_RUNTIME_DOMAIN", "lab"))        // Always passed.
    TerraformVersion  = get_env("TERRAGRUNT_CFG_BINARIES_TERRAFORM_VERSION", "1.5.1")
    TerragruntVersion = get_env("TERRAGRUNT_CFG_BINARIES_TERRAGRUNT_VERSION", "0.42.8")
    Region           = get_env("REGION", "europe-west4")
    AccountId        = get_env("LOCATION", "europe-west4")
  }
}

output "tags" {
  value       = var.tags
  description = <<EOF
The standard tags passed to this module. Ensure these tags are passed from the calling module
or the Terragrunt child configuration.
EOF
}

output "is_enabled" {
  value       = var.is_enabled
  description = "Whether the module is enabled or not"
}

output "module_input_configuration"{
  value = {
    project_id = var.project_id
    region     = var.region
    location   = var.location
    resource_friendly_identifier = var.resource_friendly_identifier
    environment = var.environment
    owner = var.owner
  }
    description = <<EOF
The input configuration passed to this module.
EOF
}

# Specific outputs for this module

/*
Global input variables.
*/
variable "is_enabled" {
  type        = bool
  description = <<EOF
  Whether this module will be created or not. It is useful, for stack-composite
modules that conditionally includes resources provided by this module..
EOF
}

variable "tags" {
  type        = map(string)
  description = "A map of tags to add to all resources."
  default     = {}
}

/*
Specific variables
*/
variable "owner" {
  type        = string
  description = <<EOF
The owner of the resource. It is used to create the resource name, so it should
be a valid username. E.g.: "john.doe" or could be a project name such as "my-project".
EOF
}

variable "resource_friendly_identifier" {
  type        = string
  description = "The name of the resource, if it's not set, it'll use the project-environment combination"
  default = null
}

variable "environment" {
  type        = string
  description = "The environment of the TSN product"
  validation {
    condition     = contains(["sandbox", "int", "stage", "prod", "master", "legacy"], var.environment )
    error_message = "The environment should be one of :sandbox, int, stage, prod, master or legacy."
  }
}

variable "location" {
  type        = string
  description = <<EOF
    The region or location where the resource will be created. It is used to
    create the resource name, so it should be a valid region or location name.
E.g.:   "EU"
EOF
  default = "europe-west4"
}

variable "region" {
  type        = string
  description = <<EOF
    The region or location where the resource will be created. It is used to
    create the resource name, so it should be a valid region or location name.
E.g.:   "eu-west1"
EOF
  default = "europe-west4"
}

variable "project_id"{
    type        = string
  description = <<EOF
The ID of the Google Cloud project where the resource will be created. It is
used to create the resource name, so it should be a valid project ID.
E.g.:  "my-project-id". For more information, refer to the documentation on
project IDs: https://cloud.google.com/resource-manager/docs/creating-managing-projects#before_you_begin
EOF
}

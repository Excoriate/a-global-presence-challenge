owner="alex"
resource_friendly_identifier="a-global-presence-function-trigger"
enable_public_http_endpoint  = true
function_name="trigger"
environment_variables = [
  {
    name = "CHALLENGE_DB_NAME"
    value = "a-global-presence-hackattic-db"
  },
  {
    name = "CHALLENGE_DOC_NAME"
    value = "challenge_doc"
  }
]

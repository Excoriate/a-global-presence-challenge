owner="alex"
resource_friendly_identifier="a-global-presence-function-shooter"
enable_public_http_endpoint  = true
function_name="shooter"
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

#location = ["europe-west4", "us-central1", "asia-east1", "asia-northeast1", "asia-southeast1", "australia-southeast1", "europe-north1", "europe-west1", "northamerica-northeast1", "southamerica-east1"]
location = ["europe-west4", "us-central1"]

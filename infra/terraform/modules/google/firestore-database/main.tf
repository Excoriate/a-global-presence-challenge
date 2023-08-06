resource "google_firestore_database" "this" {
  for_each = local.stack
  name = local.resource_name
  project = data.google_project.this[each.key].project_id
  type = "FIRESTORE_NATIVE"
  location_id = var.location
}

#resource "google_firestore_document" "this" {
#  for_each = local.stack
#  project = data.google_project.this[each.key].project_id
#  collection = "challenge_a_global_presence"
#  document_id = "alovelace"
#  database = google_firestore_database.this[each.key].name
#
#  fields = jsonencode({
#    "attempt_id" = {
#      "stringValue": uuid()
#    },
#    "token" = {
#      "stringValue": "your_token_value"
#    },
#    "attempts" = {
#      "arrayValue": {
#        "values": [
#          {
#            "mapValue": {
#              "fields": {
#                "shooterId": { "stringValue": "some_id" },
#                "isCompleted": { "booleanValue": false }
#              }
#            }
#          },
#          // Add more attempts objects if needed
#        ]
#      }
#    },
#    "countries" = {
#      "arrayValue": {
#        "values": [
#          { "stringValue": "Country1" },
#          { "stringValue": "Country2" },
#          // Add more countries if needed
#        ]
#      }
#    },
#    "IsCompleted" = {
#      "booleanValue": false
#    }
#  })
#
#  depends_on = [time_sleep.wait_60_seconds,  google_firestore_database.this]
#}

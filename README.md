# A Global presence challenge in GCP
This is an attempt to solve this [challenge](https://hackattic.com/challenges/a_global_presence) available in **Hackattic.com**. It was initially implemented in [Google Cloud Platform](https://cloud.google.com/) using [Cloud Functions](https://cloud.google.com/functions), [Firestore](https://cloud.google.com/firestore) and [Cloud Run](https://cloud.google.com/run). The code is written in [Go](https://golang.org/).

## Stack Cloud
- ☁️ [Google Cloud Platform](https://cloud.google.com/)
- ☁️ [Cloud Functions](https://cloud.google.com/functions)
- ☁️ [Firestore](https://cloud.google.com/firestore)
- ☁️ [Cloud Run](https://cloud.google.com/run)

## Development stack
- ⚙️ IaC with [Terraform](https://www.terraform.io/)
- ⚙️ [Go](https://golang.org/)
- ⚙️ [Terragrunt](https://terragrunt.gruntwork.io/) (Orchestration, and DRY)
- ⚙️ [Taskfile](https://taskfile.dev/#/) (Task runner)

---
## Configuration
### Required configuration
```bash
export PROJECT_ID="your-project-id"
export ENV=sandbox
export CHALLENGE_DB_NAME="a-global-presence-hackattic-db"
export CHALLENGE_DOC_NAME="challenge_doc"
## For testing the Go code.
## Trigger.
export IS_ENABLED=true

```
### Data layer
* The **firestore** database is called `a-global-presence-hackattic-db` and it's configured through the `common.tfvars` file located in the (`_env`) [DB stack](infra/terraform/stacks/db/hackattic).
* The actual Document that'll register the attempts and the final result is called `challenge_doc`. The [trigger](infra/terraform/stacks/cloud-function/trigger) is the responsible for creating it.

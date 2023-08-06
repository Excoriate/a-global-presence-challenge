#!/usr/bin/env bash

# Set the project id
PROJECT_ID=$1

if [ -z "$PROJECT_ID" ]; then
  echo "Project ID not set"
  exit 1
fi

# Set the details for the service account
SERVICE_ACCOUNT_ID="my-service-account"
SERVICE_ACCOUNT_DISPLAY_NAME="My Service Account"

# Function to create service account
create_service_account() {
  echo "Creating service account..."
  local service_account=$(gcloud iam service-accounts create $SERVICE_ACCOUNT_ID \
    --display-name $SERVICE_ACCOUNT_DISPLAY_NAME \
    --project $PROJECT_ID 2>&1)
    exit_status=$?
    if [ $exit_status -ne 0 ]; then
        echo "Failed to create service account: $service_account"
        exit $exit_status
    fi
}

# Function to give the service account owner access
give_owner_access() {
  echo "Giving owner access to the service account..."
  local policy=$(gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member serviceAccount:$SERVICE_ACCOUNT_ID@$PROJECT_ID.iam.gserviceaccount.com \
    --role roles/owner 2>&1)
  exit_status=$?
  if [ $exit_status -ne 0 ]; then
      echo "Failed to assign policy: $policy"
      exit $exit_status
  fi
}

# Main function to call the other cloud-function
main() {
  create_service_account
  give_owner_access
  echo "Service Account Created and Role Assigned Successfully"
}

# Call the main function
main

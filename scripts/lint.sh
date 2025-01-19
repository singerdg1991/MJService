#!/bin/bash

# Get current directory
CURR_DIR=$(dirname "$(realpath "${0}")")

# Load environment variables
. "$CURR_DIR"/base/init.sh

# Execute scripts
"$PROJECT_ROOT_DIR"/scripts/golang/vendor.sh
"$PROJECT_ROOT_DIR"/scripts/docker/docker-network.sh
"$PROJECT_ROOT_DIR"/scripts/golang/private-repo.sh
"$PROJECT_ROOT_DIR"/scripts/docker/permissions.sh

if [[ "$1" == "docker" ]]; then
  # Run the tests inside the app service
  docker compose -f "$PROJECT_ROOT_DIR"/docker-compose.local.yaml  run --rm dev_maja_service /app/scripts/prepare-lint.sh
else
  "$PROJECT_ROOT_DIR"/scripts/prepare-lint.sh
fi

# Check the exit code of the test command
if [ $? -ne 0 ]; then
  echo "Error: linter found some issues!"
  exit 1
fi
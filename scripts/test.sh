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
  docker compose -f "$PROJECT_ROOT_DIR"/docker-compose.test.yaml --profile db up -d --force-recreate --remove-orphans --build
  docker compose -f "$PROJECT_ROOT_DIR"/docker-compose.test.yaml --profile db --profile local run --rm test_maja_service_local /app/scripts/test.sh
else
  go test -v "$PROJECT_ROOT_DIR"/...
fi

# Check the exit code of the test command
if [ $? -ne 0 ]; then
  echo "Tests failed"
  exit 1
fi
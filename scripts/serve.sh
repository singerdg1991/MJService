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

"$PROJECT_ROOT_DIR"/scripts/prepare-hooks.sh

# Check if docker is used
if [[ "$1" == "docker" ]]; then
  # Check if docker is installed
  if [ -z "$(command -v docker)" ]; then
    echo "docker command not found! Please install docker!"
    exit 1
  fi

  # Check if docker compose is installed
  check_docker_compose=$( (docker compose version) 2>&1 )
  if [[ $check_docker_compose == *"is not a docker command."* ]]; then
   echo "docker compose command not found! Please install docker compose!"
   exit 1
  fi

  # Execute docker compose up
  docker compose -f \
    "$PROJECT_ROOT_DIR"/docker-compose.yaml \
    up --build --remove-orphans --force-recreate -d
else
  "$PROJECT_ROOT_DIR"/scripts/openapi.sh
	go build -o "$PROJECT_ROOT_DIR"/build/httpMaja "$PROJECT_ROOT_DIR"/cmd/http/*.go \
	 && go run "$PROJECT_ROOT_DIR"/cmd/cli/main.go start
fi

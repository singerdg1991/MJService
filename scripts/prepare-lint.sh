#!/bin/bash

# Get current directory
CURR_DIR=$(dirname "$(realpath "${0}")")

# Load environment variables
. "$CURR_DIR"/base/init.sh
"$PROJECT_ROOT_DIR"/scripts/docker/permissions.sh

# Install staticcheck if not installed
install_staticcheck() {
  if [ -z "$(command -v "${PROJECT_ROOT_DIR}"/bin/staticcheck)" ]; then
    echo "staticcheck command not found! Installing..."
    "$PROJECT_ROOT_DIR"/scripts/linter/install-staticcheck.sh --path="$PROJECT_ROOT_DIR/bin"
  fi
}

# Run staticcheck
run() {
  "$PROJECT_ROOT_DIR"/bin/staticcheck ./...
}

# install staticcheck if not installed
install_staticcheck && run

#!/bin/bash

# Get current directory
CURR_DIR=$(dirname "$(realpath "${0}")")

# Load environment variables
. "$CURR_DIR"/base/init.sh

# Install goconvey if not installed
if [ -z "$(command -v goconvey)" ]; then
  echo "goconvey command not found! Installing..."
  "$PROJECT_ROOT_DIR"/scripts/convey/install.sh --path="$PROJECT_ROOT_DIR/bin"
fi

# Run goconvey
goconvey ./... -v -race
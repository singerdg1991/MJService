#!/bin/bash

# Get current directory
CURR_DIR=$(dirname "$(realpath "${0}")")

# Load environment variables
. "$CURR_DIR"/base/init.sh

# Execute scripts
"$PROJECT_ROOT_DIR"/scripts/prepare-lint.sh

# Install task if not installed
install_task() {
  if [ -z "$(command -v "${PROJECT_ROOT_DIR}"/bin/task)" ]; then
    echo "task command not found! Installing..."
    "$PROJECT_ROOT_DIR"/scripts/watchdog/install-task.sh --path="$PROJECT_ROOT_DIR/bin"
  fi
}

# Task installation
install_task

# Run task
run() {
  "$PROJECT_ROOT_DIR"/bin/task run --watch
}

# Clean task and install again
cleanTaskAndInstall() {
  rm -rf "$PROJECT_ROOT_DIR"/bin/task &&
    rm -rf "$PROJECT_ROOT_DIR"/.task &&
    "$PROJECT_ROOT_DIR"/scripts/prepare-lint.sh &&
    install_task &&
    run
}

# Run task or clean task and install again
run || cleanTaskAndInstall

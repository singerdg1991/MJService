#!/bin/bash

# Check if go is installed
if [ -z "$(command -v go)" ]; then
  echo "go command not found! Please install go!"
  exit 1
fi

# Remove old swagger files
clean() {
  rm -rf "$PROJECT_ROOT_DIR"/public/apidocs/*.yaml
}

# Generate swagger files
generate() {
  go run "$PROJECT_ROOT_DIR"/cmd/openapi/*.go --path=./public/apidocs --internal-dir-path=./internal
}

# Clean and generate swagger files
clean && generate
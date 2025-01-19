#!/bin/bash

# Create vendor directory if not exists
if [ ! -d "$(pwd)"/vendor ]; then
  echo "vendor not found! Creating..."
  # Check go is installed
  if [ ! -x "$(command -v go)" ]; then
    echo "Go is not installed!"
  else
    go mod download &&
      go mod vendor
  fi
fi

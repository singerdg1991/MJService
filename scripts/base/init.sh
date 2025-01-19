#!/bin/bash

# Set project root directory
ROOT=$(realpath "$(dirname "${0}")"/..)

# Set environment variables
export GIT_URI=github.com
export PROJECT_ROOT_DIR=$ROOT
export SCRIPTS_DIR=$ROOT/scripts
export LOCAL_BIN_DIR=$ROOT/bin
export HEALTHCARE_DOCKER_EXTERNAL_NETWORK=hoitekBridge

# Create bin directory if not exists
mkdir -p "$PROJECT_ROOT_DIR"/bin
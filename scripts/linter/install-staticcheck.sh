#!/bin/bash

while [ $# -gt 0 ]; do
  case "$1" in
  --path=*)
    BIN_DIR="${1#*=}"
    ;;
  *)
    echo "Unknown argument: $1"
    exit 1
    ;;
  esac
  shift
done

if [ -z "$BIN_DIR" ]; then
  echo "BIN_DIR is not set"
  exit 1
fi

GOBIN=$BIN_DIR go install honnef.co/go/tools/cmd/staticcheck@latest

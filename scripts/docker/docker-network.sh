#!/bin/bash

# Check docker command exists
docker_exists=$(command -v docker)
if [ -z "$docker_exists" ]; then
  echo "Please install docker to run this command"
  exit 1
fi

# Check environment variable HEALTHCARE_DOCKER_EXTERNAL_NETWORK exists
if [ -z "$HEALTHCARE_DOCKER_EXTERNAL_NETWORK" ]; then
  echo "HEALTHCARE_DOCKER_EXTERNAL_NETWORK env is not available"
  exit 1
fi

# Check if docker network exists
docker_network_exists=$(docker network ls | grep "$HEALTHCARE_DOCKER_EXTERNAL_NETWORK")

# Create docker network if not exists
if [ -z "$docker_network_exists" ]; then
  echo "$HEALTHCARE_DOCKER_EXTERNAL_NETWORK network not found! Creating..."
  docker network create -d bridge "$HEALTHCARE_DOCKER_EXTERNAL_NETWORK"
  echo "$HEALTHCARE_DOCKER_EXTERNAL_NETWORK" " Created!"
else
  echo "$HEALTHCARE_DOCKER_EXTERNAL_NETWORK" " Already Exists!"
fi

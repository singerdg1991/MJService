#!/bin/bash

# Check .gitconfig file exists or not and create it if not exists
if [ ! -f "$PROJECT_ROOT_DIR"/.gitconfig ]; then
  echo ".gitconfig file not found, so try to create it..."

  if [[ "$GIT_CI_USERNAME" == "" ]]; then
    read -rp "Git Username: " GIT_CI_USERNAME
  fi
  if [[ "$GIT_CI_PAT" == "" ]]; then
    read -rp "Personal Access Token(PAT): " GIT_CI_PAT
  fi

  printf "[url \"https://%s:%s@%s\"]\n\tinsteadOf = %s\n" "$GIT_CI_USERNAME" "$GIT_CI_PAT" "$GIT_URI" "https://$GIT_URI" > "$PROJECT_ROOT_DIR"/.gitconfig
fi
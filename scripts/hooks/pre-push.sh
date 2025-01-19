#!/bin/sh

# Check if .PAT file exists and read it or ask for PAT
if [ -f .git-pat ]; then
  PAT=$(cat .git-pat)
else
  # Log to console
  echo "Please save your Git Personal Access Token (PAT) in .pat file in the root of the project, to update env!"
fi

# Check if .env.production file exists
if [ -f .env.production ]; then
  ENV_PROD=$(cat .env.production)
fi

# Check if .env.test file exists
if [ -f .env.test ]; then
  ENV_TEST=$(cat .env.test)
fi

# Get origin url
ORIGIN_URL=$(git config --get remote.origin.url)

# Get base url from origin url and remove ssh and port and add https
BASE_URL=$(echo "$ORIGIN_URL" | sed -e 's/ssh:\/\///g' -e 's/:.*//g' -e 's/git@/https:\/\//g')

# Update ENV_PROD in GitLab if it has value
if [ ! -z "$ENV_PROD" ]; then
  echo "Updating ENV_PROD..."
  curl --request PUT --header "PRIVATE-TOKEN: $PAT" \
       "$BASE_URL/api/v4/projects/healthcare%2Fservices%2Fmaja/variables/ENV_PROD" --form "value=$ENV_PROD"
  echo "ENV_PROD Updated!"
fi

# Update ENV_TEST in GitLab if it has value
if [ ! -z "$ENV_TEST" ]; then
  echo "Updating ENV_TEST..."
  curl --request PUT --header "PRIVATE-TOKEN: $PAT" \
       "$BASE_URL/api/v4/projects/healthcare%2Fservices%2Fmaja/variables/ENV_TEST" --form "value=$ENV_TEST"
  echo "ENV_TEST Updated!"
fi
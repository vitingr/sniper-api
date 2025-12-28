#!/bin/bash

load_env() {
  if [ -f ".env" ]; then
    export $(cat .env | grep -v '^#' | xargs)
  else
    echo "Error: .env file not found!"
    exit 1
  fi
}

load_env

if [ -z "$DATABASE_URL" ]; then
  echo "Error: DATABASE_URL not found in .env file!"
  exit 1
fi

echo "DATABASE_URL loaded: $DATABASE_URL"

case "$1" in
  up)
    migrate -database "$DATABASE_URL" -path internal/infra/database/migrations up
    ;;
  down)
    migrate -database "$DATABASE_URL" -path internal/infra/database/migrations down 1
    ;;
  create)
    read -p "Migration name: " name
    migrate create -ext sql -dir internal/infra/database/migrations -seq "$name"
    ;;
  version)
    migrate -database "$DATABASE_URL" -path internal/infra/database/migrations version
    ;;
  force)
    read -p "Version to force: " version
    migrate -database "$DATABASE_URL" -path internal/infra/database/migrations force "$version"
    ;;
  *)
    echo "Available commands: up, down, create, version, force"
    exit 1
    ;;
esac

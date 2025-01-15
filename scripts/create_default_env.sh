#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
PROJECT_ROOT_PATH=$SCRIPT_DIR"/.."
ENV_FILE_PATH=$PROJECT_ROOT_PATH"/.env"

ENV_VARS=(
    "POSTGRES_ADDRESS=postgres"
    "POSTGRES_PORT=5432"
    "POSTGRES_USER=user"
    "POSTGRES_PASSWORD=password"
    "POSTGRES_DB=order_db"

    "RABBITMQ_DEFAULT_USER=user"
    "RABBITMQ_DEFAULT_PASS=password"

    "CONFIG_PATH=config/config.yaml"
)

for VAR in "${!ENV_VARS[@]}"; do
  echo "$VAR=${ENV_VARS[$VAR]}" >> $ENV_FILE_PATH
done

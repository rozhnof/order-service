#!/bin/bash

ENV_FILE_PATH=".env"

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

echo "" > $ENV_FILE_PATH

for VAR in "${!ENV_VARS[@]}"; do
  echo "${ENV_VARS[$VAR]}" >> $ENV_FILE_PATH
done

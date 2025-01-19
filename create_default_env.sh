#!/bin/bash

ENV_FILE_PATH=".env"

ENV_VARS=(
    "PUBLISHER_ADDRESS=:8080"
    
    "POSTGRES_ADDRESS=postgres"
    "POSTGRES_PORT=5432"
    "POSTGRES_USER=user"
    "POSTGRES_PASSWORD=password"
    "POSTGRES_DB=order_db"

    "RABBITMQ_ADDRESS=rabbitmq"
    "RABBITMQ_PORT=5672"
    "RABBITMQ_DEFAULT_USER=user"
    "RABBITMQ_DEFAULT_PASS=password"

    "CONFIG_PATH=config/config.yaml"

    "MAIL_EMAIL=golang.auth.service@gmail.com"
    "MAIL_PASSWORD=jybh ayjb qosq kykn"
    "MAIL_ADDRESS=smtp.gmail.com"
    "MAIL_PORT=587"
)

echo "" > $ENV_FILE_PATH

for VAR in "${!ENV_VARS[@]}"; do
  echo "${ENV_VARS[$VAR]}" >> $ENV_FILE_PATH
done

#!/bin/sh
source ./.env
password=$USERS_POSTGRES_PASSWORD name=$USERS_POSTGRES_DATABASE user=$USERS_POSTGRES_USER port=$USERS_POSTGRES_PORT docker-compose -f stack.yml up --remove-orphans


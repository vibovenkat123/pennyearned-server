#!/bin/sh
source ./.env
db_password=$USERS_POSTGRES_PASSWORD db_name=$USERS_POSTGRES_DATABASE db_user=$USERS_POSTGRES_USER db_port=$USERS_POSTGRES_PORT redis_port=$USERS_REDIS_PORT redis_password=$USERS_REDIS_PASSWORD docker-compose -f stack.yml up --remove-orphans


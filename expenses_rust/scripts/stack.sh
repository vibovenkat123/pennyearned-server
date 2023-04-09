#!/bin/sh
source ./.env.local
password=$EXPENSES_POSTGRES_PASSWORD name=$EXPENSES_POSTGRES_DATABASE user=$EXPENSES_POSTGRES_USER port=$EXPENSES_POSTGRES_PORT docker-compose -f stack.yml up --remove-orphans


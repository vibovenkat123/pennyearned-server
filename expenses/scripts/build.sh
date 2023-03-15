#!/bin/sh
source ./.env
postgresPortCommand="main/expenses/internal/db/pkg.envPort=$EXPENSES_POSTGRES_PORT"
postgresDatabaseCommand="main/expenses/internal/db/pkg.dbName=$EXPENSES_POSTGRES_DATABASE"
postgresHostCommand="main/expenses/internal/db/pkg.dbHost=$EXPENSES_POSTGRES_HOST"
postgresPassCommand="main/expenses/internal/db/pkg.dbPass=$EXPENSES_POSTGRES_PASSWORD"
postgresUserCommand="main/expenses/internal/db/pkg.dbUser=$EXPENSES_POSTGRES_USER"
apiEnvCommand="main/expenses/internal/rest/app.env=$GO_ENV"
apiPortCommand="main/expenses/internal/rest/pkg.envPort=$EXPENSES_PORT"
if [ $GO_ENV  == "local" ]
then
    go run cmd/server/server.go
else
env GOARCH=amd64 GOOS=linux\
    go build -ldflags \
    "-X '$postgresPortCommand' -X '$postgresDatabaseCommand' -X '$postgresHostCommand' -X '$postgresPassCommand' -X '$postgresUserCommand' -X '$apiEnvCommand' -X '$apiPortCommand'"\
    -o bin/server cmd/server/server.go
fi

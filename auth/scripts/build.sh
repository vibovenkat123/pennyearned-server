#!/bin/sh
source ./.env
postgresPortCommand="main/auth/internal/db/pkg.envPort=$USERS_POSTGRES_PORT"
postgresDatabaseCommand="main/auth/internal/db/pkg.dbName=$USERS_POSTGRES_DATABASE"
postgresHostCommand="main/auth/internal/db/pkg.dbHost=$USERS_POSTGRES_HOST"
postgresPassCommand="main/auth/internal/db/pkg.dbPass=$USERS_POSTGRES_PASSWORD"
postgresUserCommand="main/auth/internal/db/pkg.dbUser=$USERS_POSTGRES_USER"
redisPortCommand="main/auth/internal/db/pkg.envRedisPort=$USERS_REDIS_PORT"
redisHostCommand="main/auth/internal/db/pkg.redisHost=$USERS_REDIS_HOST"
redisPassCommand="main/auth/internal/db/pkg.redisPass=$USERS_REDIS_PASSWORD"
apiEnvCommand="main/auth/internal/rest/app.env=$GO_ENV"
apiPortCommand="main/auth/internal/rest/pkg.envPort=$USERS_PORT"
fromEmailCommand="main/auth/internal/db/app.fromEmail=$FROM_EMAIL"
fromPasswordCommand="main/auth/internal/db/app.emailPassword=$FROM_PASSWORD"
smtpHostCommand="main/auth/internal/db/app.smtpHost=$SMTP_HOST"
smtpPortCommand="main/auth/internal/db/app.smtpPort=$SMTP_PORT"
templatePathCommand="main/auth/internal/db/app.templateFile=$TEMPLATE_PATH"
if [ $GO_ENV  == "local" ]
then
    go run cmd/server/server.go
else
    env GOARCH=amd64 GOOS=linux\
        go build -ldflags \
        "-X '$postgresPortCommand' -X '$postgresDatabaseCommand' -X '$postgresHostCommand' -X '$postgresPassCommand' -X '$postgresUserCommand' -X '$redisHostCommand' -X '$redisPassCommand' -X '$redisPortCommand' -X '$apiEnvCommand' -X '$apiPortCommand' -X '$templatePathCommand' -X '$fromEmailCommand' -X '$fromPasswordCommand' -X '$smtpHostCommand' -X '$smtpPortCommand'"\
        -o bin/server cmd/server/server.go
fi

#!/bin/bash
source ./.env
mySqlUrlCmd="main/auth/internal/db/pkg.mySqlUrl=$USERS_MYSQL_URL"
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

if [ $GO_ENV == "local" ]
then
    go run $1
else
    env GOARCH=amd64 GOOS=linux \
        go build -ldflags \
        "-X $mySqlUrlCmd -X '$redisHostCommand' -X '$redisPassCommand' -X '$redisPortCommand' -X '$apiEnvCommand' -X '$apiPortCommand' -X '$templatePathCommand' -X '$fromEmailCommand' -X '$fromPasswordCommand' -X '$smtpHostCommand' -X '$smtpPortCommand'" \
        -o $2 $1
fi


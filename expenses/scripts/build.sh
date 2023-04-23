#!/bin/sh
source ./.env
mySqlUrlCmd="main/expenses/internal/db/pkg.mySqlUrl=$EXPENSES_MYSQL_URL"
apiEnvCommand="main/expenses/internal/rest/app.env=$GO_ENV"
apiPortCommand="main/expenses/internal/rest/pkg.envPort=$EXPENSES_PORT"
if [ $GO_ENV  == "local" ]
then
    go run $1
else
env GOARCH=amd64 GOOS=linux\
    go build -ldflags \
    "-X '$mySqlUrlCmd' -X '$apiEnvCommand' -X '$apiPortCommand'"\
    -o $2 $1
fi

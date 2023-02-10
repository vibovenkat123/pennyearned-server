package dbHelpers

import (
	"os"
	"strconv"
)

var port, _ = strconv.Atoi(os.Getenv("USERS_POSTGRES_PORT"))
var redisPort, _ = strconv.Atoi(os.Getenv("USERS_REDIS_PORT"))
var DBInfo = Info{
	Host:     os.Getenv("USERS_POSTGRES_HOST"),
	Port:     port,
	User:     os.Getenv("USERS_POSTGRES_USER"),
	Password: os.Getenv("USERS_POSTGRES_PASSWORD"),
	Dbname:   os.Getenv("USERS_POSTGRES_DATABASE"),
}
var RedisInfo = Info{
    Host: os.Getenv("USERS_REDIS_HOST"),
	Port:     redisPort,
	Password: os.Getenv("USERS_REDIS_PASSWORD"),
	Dbname:   os.Getenv("USERS_REDIS_DATABASE"),
}

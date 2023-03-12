package dbHelpers

import (
	"os"
	"strconv"
)

var envPort string
var port int
var envRedisPort string
var redisPort int
var dbHost string
var dbUser string
var dbPass string
var dbName string
var redisHost string
var redisPass string
var redisName string
var DBInfo Info
var RedisInfo Info
var ConvertErr error
var RedisConvertErr error

func init() {
	if os.Getenv("GO_ENV") == "local" {
		envPort = os.Getenv("USERS_POSTGRES_PORT")
		dbUser = os.Getenv("USERS_POSTGRES_USER")
		dbHost = os.Getenv("USERS_POSTGRES_HOST")
		envRedisPort = os.Getenv("USERS_REDIS_PORT")
		dbPass = os.Getenv("USERS_POSTGRES_PASSWORD")
		dbName = os.Getenv("USERS_POSTGRES_DATABASE")
		redisHost = os.Getenv("USERS_REDIS_HOST")
		redisPass = os.Getenv("USERS_REDIS_PASSWORD")
	}
	port, ConvertErr = strconv.Atoi(envPort)
	redisPort, RedisConvertErr = strconv.Atoi(envRedisPort)
	DBInfo = Info{
		Host:     dbHost,
		Port:     port,
		User:     dbUser,
		Password: dbPass,
		Dbname:   dbName,
	}
	RedisInfo = Info{
		Host:     redisHost,
		Port:     redisPort,
		Password: redisPass,
		Dbname:   redisName,
	}
}

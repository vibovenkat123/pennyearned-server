package dbHelpers

import (
	"os"
	"strconv"
)

var envRedisPort string
var redisPort int
var mySqlUrl string
var redisHost string
var redisPass string
var redisName string
var DBInfo Info
var RedisInfo Info
var ConvertErr error
var RedisConvertErr error

func init() {
	if os.Getenv("GO_ENV") == "local" {
		mySqlUrl = os.Getenv("USERS_MYSQL_URL")
		envRedisPort = os.Getenv("USERS_REDIS_PORT")
		redisHost = os.Getenv("USERS_REDIS_HOST")
		redisPass = os.Getenv("USERS_REDIS_PASSWORD")
	}
	redisPort, RedisConvertErr = strconv.Atoi(envRedisPort)
	DBInfo = Info{
		Url: mySqlUrl,
	}
	RedisInfo = Info{
		Host:     redisHost,
		Port:     redisPort,
		Password: redisPass,
		Dbname:   redisName,
	}
}

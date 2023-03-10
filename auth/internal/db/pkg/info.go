package dbHelpers

import (
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
	port, ConvertErr = strconv.Atoi(envPort)
	redisPort, RedisConvertErr  = strconv.Atoi(envRedisPort)
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

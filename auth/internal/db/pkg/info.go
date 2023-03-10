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
func init() {
    port, _ = strconv.Atoi(envPort)
    redisPort, _ = strconv.Atoi(envRedisPort)
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

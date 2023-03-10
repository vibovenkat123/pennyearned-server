package dbHelpers

import (
	"strconv"
)

var dbHost string
var envPort string
var port int
var err error
var DBInfo Info
var dbUser string
var dbName string
var dbPass string
var ConvertErr error

func init() {
	port, ConvertErr = strconv.Atoi(envPort)
	DBInfo = Info{
		Host:     dbHost,
		Port:     port,
		User:     dbUser,
		Password: dbPass,
		Dbname:   dbName,
	}
}

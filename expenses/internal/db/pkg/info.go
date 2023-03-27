package dbHelpers

import (
	"os"
	"strconv"
)

var dbHost string
var envPort string
var port int
var DBInfo Info
var dbUser string
var dbName string
var dbPass string
var ConvertErr error

func init() {
	if os.Getenv("GO_ENV") == "local" {
		envPort = os.Getenv("EXPENSES_POSTGRES_PORT")
		dbHost = os.Getenv("EXPENSES_POSTGRES_HOST")
		dbPass = os.Getenv("EXPENSES_POSTGRES_PASSWORD")
		dbUser = os.Getenv("EXPENSES_POSTGRES_USER")
		dbName = os.Getenv("EXPENSES_POSTGRES_DATABASE")
	}
	port, ConvertErr = strconv.Atoi(envPort)
	DBInfo = Info{
		Host:     dbHost,
		Port:     port,
		User:     dbUser,
		Password: dbPass,
		Dbname:   dbName,
	}
}

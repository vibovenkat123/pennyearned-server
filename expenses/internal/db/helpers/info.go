package dbHelpers

import (
	"os"
	"strconv"
)

var port, _ = strconv.Atoi(os.Getenv("EXPENSES_POSTGRES_PORT"))
var DBInfo = Info{
	Host:     os.Getenv("EXPENSES_POSTGRES_HOST"),
	Port:     port,
	User:     os.Getenv("EXPENSES_POSTGRES_USER"),
	Password: os.Getenv("EXPENSES_POSTGRES_PASSWORD"),
	Dbname:   os.Getenv("EXPENSES_POSTGRES_DATABASE"),
}

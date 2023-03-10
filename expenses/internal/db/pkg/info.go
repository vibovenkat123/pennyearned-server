package dbHelpers

import (
	"strconv"
    "log"
)
var dbHost string
var envPort string
var port int
var err error
var DBInfo Info
var dbUser string
var dbName string
var dbPass string
func init() {
    port, err = strconv.Atoi(envPort) 
    if err != nil {
        log.Fatalln(err)
    }
    DBInfo = Info{
        Host:     dbHost,
        Port:     port,
        User:     dbUser,
        Password: dbPass,
        Dbname:   dbName,
    }
}

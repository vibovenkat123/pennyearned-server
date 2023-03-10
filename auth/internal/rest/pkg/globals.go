package apiHelpers

import (
    "strconv"
    "log"
)

type UserIDRes struct {
	ID string `json:"id"`
}

var Port int
var err error
var envPort string
func init() {
    Port, err = strconv.Atoi(envPort)
    if err != nil {
        log.Fatalln(err)
    }
}

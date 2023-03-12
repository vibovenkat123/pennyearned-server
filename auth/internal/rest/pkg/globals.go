package apiHelpers

import (
	"os"
	"strconv"
)

type UserIDRes struct {
	ID string `json:"id"`
}

var Port int
var ConvertErr error
var envPort string
var Local bool

func init() {
	Local = os.Getenv("GO_ENV") == "local"
	if Local {
		envPort = os.Getenv("USERS_PORT")
	}
	Port, ConvertErr = strconv.Atoi(envPort)
}

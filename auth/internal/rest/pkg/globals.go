package apiHelpers

import (
	"os"
	"strconv"
)

type UserIDRes struct {
	ID string `json:"id"`
}

type UserAccessRes struct {
	AccessToken string `json:"access_token"`
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

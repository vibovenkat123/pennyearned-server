package apiHelpers

import (
	"os"
	"strconv"
)

type SignInData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpVerifyData struct {
	Name     string `json:"name,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type VerifyData struct {
	Email string `json:"email"`
}

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

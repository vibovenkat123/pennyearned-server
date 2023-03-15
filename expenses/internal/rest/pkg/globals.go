package apiHelpers

import (
	"os"
	"strconv"
)

type UpdateExpenseData struct {
	Name  string `json:"name"`
	Spent int    `json:"spent"`
}
type NewExpenseData struct {
	OwnerID string `json:"ownerid"`
	Name    string `json:"name"`
	Spent   int    `json:"spent"`
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

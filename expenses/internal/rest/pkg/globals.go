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
// the type to envelope json in a key
type Envelope map[string]any
var Port int
var ConvertErr error
var envPort string
var Local bool

// the default key the envelope uses
var topKey = "expense"


func init() {
	Local = os.Getenv("GO_ENV") == "local"
	if Local {
		envPort = os.Getenv("EXPENSES_PORT")
	}
	Port, ConvertErr = strconv.Atoi(envPort)
}

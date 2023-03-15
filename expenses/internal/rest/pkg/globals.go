package apiHelpers

import (
	"os"
	"strconv"

	"go.uber.org/zap"
)

type UpdateExpenseData struct {
	Name  string `json:"name"`
	Spent int    `json:"spent"`
}

type Application struct {
	Log *zap.Logger
	Conf Config
}


type Config struct {
	Port int
}

var App *Application = &Application{}

func SetLogger(logger *zap.Logger) {
	App.Log = logger
}

type NewExpenseData struct {
	OwnerID string `json:"ownerid"`
	Name    string `json:"name"`
	Spent   int    `json:"spent"`
}

// the type to envelope json in a key
type Envelope map[string]any

var ConvertErr error
var Local bool
var err error
// the default key the envelope uses
var topKey = "expense"

func init() {
	var envPort string
	Local = os.Getenv("GO_ENV") == "local"
	if Local {
		envPort = os.Getenv("EXPENSES_PORT")
	}
	port, err := strconv.Atoi(envPort)
	if err != nil {
		App.Log.Error("Failed to convert port to int",
			zap.Error(err),
		)
	}
	App.Conf.Port = port

}

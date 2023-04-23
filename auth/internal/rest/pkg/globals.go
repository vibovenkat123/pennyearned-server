package apiHelpers

import (
	"os"
	"strconv"

	"go.uber.org/zap"
)

type Application struct {
	Log  *zap.Logger
	Conf Config
}

type Config struct {
	Port int
}

var App *Application = &Application{}

func SetLogger(logger *zap.Logger) {
	App.Log = logger
}

type SignInData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpVerifyData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Make the json pretty printed/tabbed
type Envelope map[string]any

type VerifyData struct {
	Email string `json:"email"`
}

type UserAccessRes struct {
	AccessToken string `json:"access_token"`
}

// the top key that the envelope uses by default
var topKey = "user"
var Local bool

func Initialize() {
	var envPort string
	Local = os.Getenv("GO_ENV") == "local"
	if Local {
		envPort = os.Getenv("USERS_PORT")
	}
	port, err := strconv.Atoi(envPort)
	if err != nil {
		App.Log.Error("Failed to convert port to int",
			zap.Error(err),
		)
	}
	App.Conf.Port = port
}

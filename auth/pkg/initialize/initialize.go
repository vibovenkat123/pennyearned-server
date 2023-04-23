package initialize

import (
	database "main/auth/internal/db/app"
	api "main/auth/internal/rest/app"
	apiGlobals "main/auth/internal/rest/pkg"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func Initialize() {
	Logger, _ = zap.NewProduction()
	defer Logger.Sync()
	// init the logger
	database.InitializeLogger(Logger)
	// connect to database
	db, _ := database.Connect(Logger)
	// check if we successfully connected
	err := db.Ping()
	if err != nil {
		Logger.Panic("Cannot ping the database",
			zap.Error(err),
		)
	}
	// expose endpoints
	apiGlobals.SetLogger(Logger)
	apiGlobals.Initialize()
	api.StartAPI()
}

package initialize

import (
	"go.uber.org/zap"
	database "main/expenses/internal/db/app"
	api "main/expenses/internal/rest/app"
	apiGlobals "main/expenses/internal/rest/pkg"
)

var Logger *zap.Logger

func Initialize() {
	Logger, _ = zap.NewProduction()
	defer Logger.Sync()
	// connect to database
	db, err := database.Connect(Logger)
	if err != nil {
		Logger.Error("Error while connecting to database",
			zap.Error(err),
		)
	}
	// check if we successfully connected
	err = db.Ping()
	if err != nil {
		Logger.Error("Error while pinging database",
			zap.Error(err),
		)
	}
	database.InitializeLogger(Logger)
	// migrate (add columns and tables)
	apiGlobals.SetLogger(Logger)
	apiGlobals.Initialize()
	// expose endpoints
	api.StartAPI()
}

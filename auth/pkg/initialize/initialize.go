package initialize

import (
	"go.uber.org/zap"
	database "main/auth/internal/db/app"
	api "main/auth/internal/rest/app"
	apiGlobals "main/auth/internal/rest/pkg"
)

var Logger *zap.Logger

func init() {
	Logger, _ = zap.NewProduction()
	defer Logger.Sync()
}
func Initialize() {
	// connect to database
	db, _ := database.Connect(Logger)
	// check if we successfully connected
	err := db.Ping()
	if err != nil {
		Logger.Panic("Cannot ping the database",
			zap.Error(err),
		)
	}
	database.InitializeLogger(Logger)
	// expose endpoints
	apiGlobals.SetLogger(Logger)
	api.StartAPI()
}

package initialize

import (
	"go.uber.org/zap"
	database "main/expenses/internal/db/app"
	api "main/expenses/internal/rest/app"
)

var Logger *zap.Logger

func init() {
	Logger, _ = zap.NewProduction()
	defer Logger.Sync()
}
func Initialize() {
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
	// migrate (add columns and tables)
	database.Migrate()
	// WARNING: THE FOLLOWING LINE WILL
	// DESTROY: THE DATABASE
	// database.ResetToSchema()
	// expose endpoints
	api.StartAPI(Logger)
}

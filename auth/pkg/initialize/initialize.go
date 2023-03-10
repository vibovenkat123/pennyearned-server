package initialize

import (
	"go.uber.org/zap"
	database "main/auth/internal/db/app"
	api "main/auth/internal/rest/app"
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
	// migrate (add columns and tables)
	database.Migrate()
	// WARNING: THE FOLLOWING LINE WILL DESTROY THE DATABASE
	database.ResetToSchema()
	// expose endpoints
	api.Expose(Logger)
}

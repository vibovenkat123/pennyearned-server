package main

import (
	dbHelpers "main/auth/internal/db/helpers"
	database "main/auth/internal/db/services"
	api "main/auth/internal/rest/services"
)
func main() {
	// connect to database
	db, dbErr, _, redisErr := database.Connect()
	if dbErr != nil {
		panic(dbErr)
	}
    if redisErr != nil {
        panic(redisErr)
    }
	// check if we successfully connected
	dbErr = db.Ping()
	if dbErr != nil {
		panic(dbErr)
	}
	// migrate (add columns and tables)
	dbHelpers.Migrate()
	// WARNING: THE FOLLOWING LINE WILL
	// DESTROY: THE DATABASE
	// dbHelpers.ResetToSchema()
	// expose endpoints
	api.Expose()
}

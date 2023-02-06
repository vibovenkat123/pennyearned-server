package main

import (
	dbHelpers "main/expenses/internal/db/helpers"
	database "main/expenses/internal/db/services"
    api "main/expenses/internal/rest/services"
)

func main() {
	// connect to database
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	// check if we successfully connected
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	// migrate (add columns and tables)
	dbHelpers.Migrate()
	// WARNING: THE FOLLOWING LINE WILL
	// DESTROY: THE DATABASE
	// dbHelpers.ResetToSchema()
	// expose endpoints
    api.Expose()
}

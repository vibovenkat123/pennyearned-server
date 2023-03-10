package main

import (
	database "main/expenses/internal/db/app"
	api "main/expenses/internal/rest/app"
    "log"
)

func main() {
	// connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	// check if we successfully connected
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	// migrate (add columns and tables)
	database.Migrate()
	// WARNING: THE FOLLOWING LINE WILL
	// DESTROY: THE DATABASE
	// dbHelpers.ResetToSchema()
	// expose endpoints
	api.Expose()
}

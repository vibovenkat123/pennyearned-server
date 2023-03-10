package main

import (
    database "main/auth/internal/db/app"
	api "main/auth/internal/rest/app"
    "log"
)

func main() {
	// connect to database
	db, _:= database.Connect()
    // check if we successfully connected
    err := db.Ping()
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

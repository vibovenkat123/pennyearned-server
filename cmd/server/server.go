package main

import (
//	dbHelpers "github.com/vibovenkat123/pennyearned-server/internal/db/helpers"
	database "github.com/vibovenkat123/pennyearned-server/internal/db/services"
	api "github.com/vibovenkat123/pennyearned-server/internal/rest/services"
)

func main() {
	// connect to database
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// WARNING: enable if you want to reset to the schema
	// DANGER: ENABLING THE COMMAND WILL DELETE ALL THE TABLES
//	dbHelpers.ResetToSchema(db)
	api.Expose()
}

package main

import (
	dbHelpers "github.com/vibovenkat123/pennyearned-server/internal/db/helpers"
	database "github.com/vibovenkat123/pennyearned-server/internal/db/services"
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
	users := []dbHelpers.User{}
	ownerid := "8ff0c79d-adeb-482a-9bca-dd7687f5cac3"
	dbHelpers.Db.Select(&users, "SELECT * FROM users where id=$1", ownerid)
	// migrate (add columns and tables0
	dbHelpers.Migrate()
	// WARNING: THE FOLLOWING LINE WILL
	// DESTROY: THE DATABASE
	//    dbHelpers.ResetToSchema()
	// expose endpoints
}

package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	helpers "github.com/vibovenkat123/pennyearned-server/internal/db/helpers"
)

func Connect() (helpers.DatabaseType, error) {
	var err error
	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		helpers.DbInfo.Host, helpers.DbInfo.Port, helpers.DbInfo.User, helpers.DbInfo.Password, helpers.DbInfo.Dbname)
	fmt.Println("Attempting to connect...")
	helpers.Db, err = sqlx.Open("postgres", postgresqlDbInfo)
	if err != nil {
		log.Fatalln(err)
	}
	err = helpers.Db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected!!")
	return helpers.Db, err
}

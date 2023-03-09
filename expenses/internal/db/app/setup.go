package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	helpers "main/expenses/internal/db/pkg"
)

func Connect() (helpers.DatabaseType, error) {
	var err error
	expensesDBInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		helpers.DBInfo.Host, helpers.DBInfo.Port, helpers.DBInfo.User, helpers.DBInfo.Password, helpers.DBInfo.Dbname)
	fmt.Println("Attempting to connect...")
	helpers.DB, err = sqlx.Open("postgres", expensesDBInfo)
	if err != nil {
		log.Fatalln(err)
	}
	err = helpers.DB.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected!!")
	return helpers.DB, err
}

package db

import (
	"fmt"
	"log"
    "github.com/redis/go-redis/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	helpers "main/expenses/internal/db/helpers"
)

func Connect() (helpers.DatabaseType, error) {
	var err error
    helpers.RDB = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        Password: "Ppo39UTaN5pJnWrDEQtCHf9Mn8ycSnJKG",
        DB: 0,
    })
	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		helpers.DbInfo.Host, helpers.DbInfo.Port, helpers.DbInfo.User, helpers.DbInfo.Password, helpers.DbInfo.Dbname)
	fmt.Println("Attempting to connect...")
	helpers.DB, err = sqlx.Open("postgres", postgresqlDbInfo)
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

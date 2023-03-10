package db

import (
	"context"
	"fmt"
	"log"
    helpers "main/auth/internal/db/pkg"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func Connect() (helpers.DatabaseType, helpers.RedisType) {
    rdb, redisErr := ConnectRedis()
    if redisErr != nil {
        log.Fatalln(redisErr)
    }
	db, dbErr := ConnectPostgres()
    if dbErr != nil {
        log.Fatalln(dbErr)
    }
	return db, rdb 
}
func ConnectRedis() (helpers.RedisType, error) {
    var dbNumber int
    var err error
    if len(helpers.RedisInfo.Dbname) == 0{ 
        dbNumber = 0
    } else {
        dbNumber, err  = strconv.Atoi(helpers.RedisInfo.Dbname)
        if err != nil {
            log.Fatalln(err)
        }
    }
	helpers.RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", helpers.RedisInfo.Host, helpers.RedisInfo.Port),
		Password: helpers.RedisInfo.Password,
		DB:       dbNumber,
	})
	err = helpers.RDB.Ping(ctx).Err()
	return helpers.RDB, err
}
func ConnectPostgres() (helpers.DatabaseType, error) {
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

package db

import (
	"context"
	"fmt"
	"log"
	helpers "main/auth/internal/db/helpers"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func Connect() (helpers.DatabaseType, error, helpers.RedisType, error) {
	db, err := ConnectPostgres()
	rdb, redisErr := ConnectRedis()
	return db, err, rdb, redisErr
}
func ConnectRedis() (helpers.RedisType, error) {
	dbNumber, _ := strconv.Atoi(helpers.RedisInfo.Dbname)
	helpers.RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", helpers.RedisInfo.Host, helpers.RedisInfo.Port),
		Password: helpers.RedisInfo.Password,
		DB:       dbNumber,
	})
	err := helpers.RDB.Ping(ctx).Err()
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

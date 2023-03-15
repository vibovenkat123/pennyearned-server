package db

import (
	"context"
	"fmt"
	helpers "main/auth/internal/db/pkg"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

func Connect(log *zap.Logger) (helpers.DatabaseType, helpers.RedisType) {
	if helpers.ConvertErr != nil {
		log.Error("Failed to convert port env to integer for postgres",
			zap.Error(helpers.ConvertErr),
		)
	}
	if helpers.RedisConvertErr != nil {
		log.Error("Failed to convert port env to integer for redis",
			zap.Error(helpers.RedisConvertErr),
		)
	}
	rdb, redisErr := ConnectRedis()
	if redisErr != nil {
		log.Error("Failed to connect to redis",
			zap.Error(redisErr),
		)
	}
	db, dbErr := ConnectPostgres()
	if dbErr != nil {
		log.Error("Failed to connect to postgres database",
			zap.Error(dbErr),
		)
	}
	return *db, *rdb
}
func ConnectRedis() (*helpers.RedisType, error) {
	var dbNumber int
	var err error
	if len(helpers.RedisInfo.Dbname) == 0 {
		dbNumber = 0
	} else {
		dbNumber, err = strconv.Atoi(helpers.RedisInfo.Dbname)
		if err != nil {
			return nil, err
		}
	}
	helpers.RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", helpers.RedisInfo.Host, helpers.RedisInfo.Port),
		Password: helpers.RedisInfo.Password,
		DB:       dbNumber,
	})
	err = helpers.RDB.Ping(ctx).Err()
	return &helpers.RDB, err
}
func ConnectPostgres() (*helpers.DatabaseType, error) {
	var err error
	usersDBInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		helpers.DBInfo.Host, helpers.DBInfo.Port, helpers.DBInfo.User, helpers.DBInfo.Password, helpers.DBInfo.Dbname)
	fmt.Println("Attempting to connect...")
	helpers.DB, err = sqlx.Open("postgres", usersDBInfo)
	if err != nil {
		return nil, err
	}
	err = helpers.DB.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected!!")
	return &helpers.DB, err
}

package db

import (
	"context"
	"fmt"
	helpers "main/auth/internal/db/pkg"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

func Connect(log *zap.Logger) (helpers.DatabaseType, helpers.RedisType) {
	if helpers.RedisConvertErr != nil {
		log.Error("Failed to convert port env to integer for redis",
			zap.Error(helpers.RedisConvertErr),
		)
	}
	log.Info("Attempting to connect...")
	rdb, redisErr := ConnectRedis()
	if redisErr != nil {
		log.Panic("Failed to connect to redis",
			zap.Error(redisErr),
		)
	}
	db, dbErr := ConnectMySQL()
	if dbErr != nil {
		log.Panic("Failed to connect to mysql database",
			zap.Error(dbErr),
		)
	}
	return *db, *rdb
}
func ConnectRedis() (*helpers.RedisType, error) {
	var dbNumber int
	var err error
	log.Info("Getting db name")
	if len(helpers.RedisInfo.Dbname) == 0 {
		dbNumber = 0
	} else {
		dbNumber, err = strconv.Atoi(helpers.RedisInfo.Dbname)
		if err != nil {
			return nil, err
		}
	}
	log.Info("Got Name",
		zap.Int("Name", dbNumber),
	)
	addr := fmt.Sprintf("%v:%v", helpers.RedisInfo.Host, helpers.RedisInfo.Port)
	log.Info("Connecting to redis",
		zap.String("Addr", addr),
		zap.String("Password", helpers.RedisInfo.Password),
		zap.Int("Name", dbNumber),
	)
	helpers.RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: helpers.RedisInfo.Password,
		DB:       dbNumber,
		Username: "default",
	})
	log.Info("Pinging redis")
	err = helpers.RDB.Ping(ctx).Err()
	return &helpers.RDB, err
}
func ConnectMySQL() (*helpers.DatabaseType, error) {
	var err error
	log.Info("Connecting to mysql database",
		zap.String("Url", helpers.DBInfo.Url),
	)
	helpers.DB, err = sqlx.Open("mysql", helpers.DBInfo.Url)
	if err != nil {
		return nil, err
	}
	log.Info("Pinging database")
	err = helpers.DB.Ping()
	return &helpers.DB, err
}

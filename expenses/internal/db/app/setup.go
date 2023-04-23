package db

import (
	"fmt"
	helpers "main/expenses/internal/db/pkg"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func Connect(log *zap.Logger) (helpers.DatabaseType, error) {
	var err error
	log.Info("Attempting to connect...",
		zap.String("URL", helpers.DBInfo.Url),
	)
	helpers.DB, err = sqlx.Open("mysql", helpers.DBInfo.Url)
	if err != nil {
		return nil, err
	}
	log.Info("Pinging database")
	err = helpers.DB.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected!!")
	return helpers.DB, err
}

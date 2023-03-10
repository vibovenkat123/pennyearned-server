package db

import (
	"fmt"
	helpers "main/expenses/internal/db/pkg"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func Connect(log *zap.Logger) (helpers.DatabaseType, error) {
    if helpers.ConvertErr != nil {
        log.Error("The port variable is not a valid int",
            zap.Error(helpers.ConvertErr),
        )
    }
	var err error
	expensesDBInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		helpers.DBInfo.Host, helpers.DBInfo.Port, helpers.DBInfo.User, helpers.DBInfo.Password, helpers.DBInfo.Dbname)
	fmt.Println("Attempting to connect...")
	helpers.DB, err = sqlx.Open("postgres", expensesDBInfo)
	if err != nil {
        return nil, err
	}
	err = helpers.DB.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected!!")
	return helpers.DB, err
}

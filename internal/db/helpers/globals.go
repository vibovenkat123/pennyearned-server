package dbHelpers

import (
	"github.com/jmoiron/sqlx"
)

type Expenses struct {
	ID      string `db:"id"`
	OwnerId string `db:"owner_id"`
	Name    string `db:"name"`
	Spent   int    `db:"spent"`
}
type User struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

var Db *sqlx.DB

type DatabaseType = *sqlx.DB
type Info struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

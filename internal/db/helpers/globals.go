package dbHelpers

import (
	"github.com/jmoiron/sqlx"
)

type Schema struct {
	create string
	drop   string
	alter  string
}
type Expense struct {
	ID      string `db:"id"`
	OwnerID string `db:"owner_id"`
	Name    string `db:"name"`
	Spent   int    `db:"spent"`
    Date string `db:"date_created"`
}
type User struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
    Date string `db:"date_created"`
}

type DatabaseType = *sqlx.DB

var Db DatabaseType

type Info struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

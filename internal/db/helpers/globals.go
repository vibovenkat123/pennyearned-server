package dbHelpers

import (
	"github.com/jmoiron/sqlx"
)
type Schema struct {
    create string
    drop string
    alter string
}
type Expense struct {
	ID      string `db:"id"`
	OwnerID string `db:"owner_id"`
	Name    string `db:"name"`
	Spent   int    `db:"spent"`
}
type User struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
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

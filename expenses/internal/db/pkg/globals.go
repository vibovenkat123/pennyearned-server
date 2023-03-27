package dbHelpers

// globals
import (
	"errors"
	"github.com/jmoiron/sqlx"
)

type IDResponse struct {
	ID string `json:"id"`
}

var (
	ErrEmailInvalid     = errors.New("Email is invalid")
	ErrExpensesNotFound = errors.New("Expenses not found")
	ErrExpenseNotFound  = errors.New("Expense not found")
	ErrInvalidFormat    = errors.New("Invalid format")
)

type Schema struct {
	Create string
	Drop   string
	Alter  string
}
type Expense struct {
	ID          string `db:"id" json:"id"`
	OwnerID     string `db:"owner_id" json:"owner_id"`
	Name        string `db:"name" json:"name"`
	DateCreated string `db:"date_created" json:"date_created"`
	DateUpdated string `db:"date_updated" json:"date_updated"`
	Spent       int    `db:"spent" json:"spent"`
}
type DatabaseType = *sqlx.DB

var DB DatabaseType

type Info struct {
	Host     string
	User     string
	Password string
	Dbname   string
	Port     int
}

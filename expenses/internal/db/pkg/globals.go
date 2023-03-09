package dbHelpers

// globals
import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Response struct {
	ID string `json:"id"`
}

var (
	minNameLength     = 1
	maxNameLength     = 20
	minUsernameLength = 2
	maxUsernameLength = 30
	minPasswordLength = 8
	maxPasswordLength = 20
)
var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
	ErrPassNotMatch        = errors.New("Password is invalid")
	ErrEmailNotFound       = errors.New("Email not found")
	ErrAlreadyFound        = errors.New("Email or Username already found")
	ErrUsernameTooShort    = errors.New(fmt.Sprintf("Username must be bigger than %v character(s)", minUsernameLength))
	ErrUsernameTooLong     = errors.New(fmt.Sprintf("Username must be smaller than %v character(s)", maxUsernameLength))
	ErrPasswordTooShort    = errors.New(fmt.Sprintf("Password must be bigger than %v character(s)", minPasswordLength))
	ErrPasswordTooLong     = errors.New(fmt.Sprintf("Password must be smaller than %v character(s)", maxPasswordLength))
	ErrNameTooShort        = errors.New(fmt.Sprintf("Name must be bigger than %v character(s)", minNameLength))
	ErrNameTooLong         = errors.New(fmt.Sprintf("Name must be smaller than %v character(s)", maxNameLength))
	ErrEmailInvalid        = errors.New("Email is invalid")
	ErrExpensesNotFound    = errors.New("Expenses not found")
	ErrExpenseNotFound     = errors.New("Expense not found")
	ErrInvalidFormat       = errors.New("Invalid format")
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

type Schema struct {
	Create string
	Drop   string
	Alter  string
}
type Expense struct {
	ID          string `db:"id"`
	OwnerID     string `db:"owner_id"`
	Name        string `db:"name"`
	Spent       int    `db:"spent"`
	DateCreated string `db:"date_created"`
	DateUpdated string `db:"date_updated"`
}
type User struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Email       string `db:"email"`
	Username    string `db:"username"`
	Password    string `db:"password"`
	DateCreated string `db:"date_created"`
	DateUpdated string `db:"date_updated"`
}
type DatabaseType = *sqlx.DB

var DB DatabaseType
var P = &params{
	iterations:  3,
	parallelism: 2,
	saltLength:  16,
	memory:      64 * 1024,
	keyLength:   32,
}

type Info struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

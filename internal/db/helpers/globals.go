package dbHelpers

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
	ErrPassNotMatch        = errors.New("Password is invalid")
	ErrEmailNotFound       = errors.New("Email is invalid")
	ErrUsernameTooShort = errors.New("Username is too short")
	ErrUsernameTooLong = errors.New("Username is too long")
	ErrPasswordTooShort =  errors.New("Username is too short")
	ErrPasswordTooLong=  errors.New("Username is too long")
    ErrNameTooShort = errors.New("Name is too short")
    ErrNameTooLong = errors.New("Name is too long")
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

type Schema struct {
	create string
	drop   string
	alter  string
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

var Db DatabaseType
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

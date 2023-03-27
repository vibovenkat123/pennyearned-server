package dbHelpers

// globals
import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type IDResponse struct {
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
	ErrInvalidCode         = errors.New("The code you entered is not found")
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
)

type Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type User struct {
	ID          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name,omitempty"`
	Email       string `db:"email" json:"email"`
	Username    string `db:"username" json:"username"`
	Password    string `db:"password" json:"password"`
	DateCreated string `db:"date_created" json:"date_created"`
	DateUpdated string `db:"date_updated" json:"date_updated"`
}
type DatabaseType = *sqlx.DB
type RedisType = *redis.Client

var DB DatabaseType
var RDB RedisType
var P = &Params{
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	Memory:      64 * 1024,
	KeyLength:   32,
}

type Info struct {
	Host     string
	User     string
	Password string
	Dbname   string
	Port     int
}

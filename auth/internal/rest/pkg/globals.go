package apiHelpers

import (
	"os"
	"strconv"
)

type UserIDRes struct {
	ID string `json:"id"`
}

var Port, Err = strconv.Atoi(os.Getenv("USERS_PORT"))

package apiHelpers

import (
	"os"
	"strconv"
)
var Port, Err = strconv.Atoi(os.Getenv("EXPENSES_PORT"))

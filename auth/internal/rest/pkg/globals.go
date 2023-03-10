package apiHelpers

import (
	"strconv"
)

type UserIDRes struct {
	ID string `json:"id"`
}

var Port int
var ConvertErr error
var envPort string

func init() {
	Port, ConvertErr = strconv.Atoi(envPort)
}

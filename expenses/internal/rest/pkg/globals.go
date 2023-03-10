package apiHelpers

import (
	"strconv"
)

var Port int
var ConvertErr error
var envPort string

func init() {
	Port, ConvertErr = strconv.Atoi(envPort)
}

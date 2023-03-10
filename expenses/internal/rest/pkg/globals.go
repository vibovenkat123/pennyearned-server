package apiHelpers

import (
	"strconv"
    "log"
)
var Port int
var err error
var envPort string
func init() {
   Port, err = strconv.Atoi(envPort)
   if err != nil {
        log.Fatalln(err)
   }
}

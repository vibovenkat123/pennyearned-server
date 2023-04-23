package dbHelpers

import (
	"os"
)

var mySqlUrl string
var DBInfo Info

func init() {
	if os.Getenv("GO_ENV") == "local" {
		mySqlUrl = os.Getenv("EXPENSES_MYSQL_URL")
	}
	DBInfo = Info{
		Url: mySqlUrl,
	}
}

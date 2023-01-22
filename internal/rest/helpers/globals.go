package apiHelpers

import (
	dbHelpers "github.com/vibovenkat123/pennyearned-server/internal/db/helpers"
)

type DatabaseType = dbHelpers.DatabaseType
type AppConfig struct {
	Port int
	Env  string
}

var Conf AppConfig

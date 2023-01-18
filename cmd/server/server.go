package main

import (
    setup "github.com/vibovenkat123/pennyearned-server/pkg/setup"
)

func main() {
    // connect to database
    db, err := setup.Connect()
    if err != nil {
        panic(err)
    }
    // WARNING: enable if you want to reset to the scheme
    // DANGER: ENABLING THE COMMAND WILL DELETE ALL THE TABLES
//    setup.ResetToSchema(db)
}

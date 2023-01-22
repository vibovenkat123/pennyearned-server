package dbHelpers

import (
	"fmt"
	"log"
	"strings"
)

func ResetToSchema(db DatabaseType) {
	fmt.Println("Resetting...")
	ExecMultiple(db, defaultSchema.drop)
	db.MustExec(defaultSchema.create)
	ExecMultiple(db, defaultSchema.alter)
	fmt.Println("Resetted!!")
}
func GetExpensesByOwnerId(ownerid string) []Expenses {
	expenses := []Expenses{}
	Db.Select(&expenses, "SELECT * FROM expenses where owner_id=$1", ownerid)
	return expenses
}
func GetOwnerById(ownerid string) User {
	user := User{}
	Db.Select(&user, "SELECT * FROM users where _id=$1", ownerid)
	return user
}
func ExecMultiple(e DatabaseType, query string) {
	statements := strings.Split(query, "\n")
	if len(strings.Trim(statements[len(statements)-1], " \n\t\r")) == 0 {
		statements = statements[:len(statements)-1]
	}
	for _, s := range statements {
		_, err := e.Exec(s)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

package dbHelpers

import (
	"fmt"
	"log"
	"strings"
)
func Migrate() {
    fmt.Println("Migrating...")
    Db.MustExec(defaultSchema.create)
    ExecMultiple(Db, defaultSchema.alter)
    fmt.Println("Migrated!!")
}
func ResetToSchema() {
	fmt.Println("Resetting...")
	ExecMultiple(Db, defaultSchema.drop)
	Db.MustExec(defaultSchema.create)
	ExecMultiple(Db, defaultSchema.alter)
	fmt.Println("Resetted!!")
}
func GetExpenseById(expenseId string) Expense {
	expense := Expense{}
	Db.Get(&expense, "SELECT * FROM expenses where id=$1", expenseId)
	return expense
}
func GetExpensesByOwnerId(ownerid string) []Expense {
	expenses := []Expense{}
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

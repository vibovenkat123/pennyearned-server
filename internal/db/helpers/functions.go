package dbHelpers

import (
	"fmt"
	"log"
	"strings"
    "github.com/google/uuid"
//    "reflect"
)
func GetAllUsers() []User {
    var arr []User
	Db.Select(&arr, "SELECT * FROM users ORDER BY name ASC")
    return arr
}
func GetAllExpenses() []Expense {
    var arr []Expense
	Db.Select(&arr, "SELECT * FROM expenses ORDER BY name ASC")
    return arr
}
func NewExpense(ownerid string, name string, spent int) Expense{
    tx := Db.MustBegin()
    id  := uuid.New()
    tx.MustExec("INSERT INTO expenses (owner_id, name, spent, id) VALUES ($1, $2, $3, $4)", ownerid, name, spent, id.String())
    tx.Commit()
    return GetExpenseById(id.String())
}
func DeleteExpense(id string) Expense{
    expense := GetExpenseById(id)
    if expense.ID == "" {
        expense.ID = "not found"
        expense.Name = "not found"
        expense.OwnerID= "not found"
        expense.Spent = 42
        return expense
    }
    Db.MustExec(`DELETE FROM expenses WHERE id=$1`, id)
    return expense
}
func UpdateExpense(expense Expense) Expense {
    Db.MustExec(`UPDATE expenses SET owner_id=$1, name=$2, spent=$3 WHERE id=$4`, expense.OwnerID, expense.Name, expense.Spent, expense.ID)
    return GetExpenseById(expense.ID)
}
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

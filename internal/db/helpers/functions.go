package dbHelpers

import (
	"fmt"
	"log"
	"strings"
    "github.com/google/uuid"
//    "reflect"
)
func GetAllExpenses() []Expense {
    var expenses []Expense Db.Select(&expenses, "SELECT * FROM expenses ORDER BY date_created ASC")
    return expenses
}
func MostExpensiveExpense(ownerid string) Expense{
    expensesData := GetExpensesByOwnerId(ownerid)
    mostExpensive := expensesData[0].Spent
    id := expensesData[0].ID
    for _, i := range expensesData[1:] {
        if i.Spent > mostExpensive{
            mostExpensive = i.Spent
            id = i.ID
        }
    }
    return GetExpenseById(id)
}
func LeastExpensiveExpense(ownerid string) Expense{
    expensesData := GetExpensesByOwnerId(ownerid)
    leastExpensive := expensesData[0].Spent
    id := expensesData[0].ID
    for _, i := range expensesData[1:] {
        if i.Spent < leastExpensive{
            leastExpensive = i.Spent
            id = i.ID
        }
    }
    return GetExpenseById(id)
}
func ExpensesLowerThan(spent int, ownerid string) []Expense {
    var expenses []Expense
    expensesData := GetExpensesByOwnerId(ownerid)
    for _, i := range expensesData {
        if i.Spent < spent{
            expenses = append(expenses, i)
        }
    }
    return expenses
}
func ExpensesHigherThan(spent int, ownerid string) []Expense {
    var expenses []Expense
    expensesData := GetExpensesByOwnerId(ownerid)
    for _, i := range expensesData {
        if i.Spent > spent{
            expenses = append(expenses, i)
        }
    }
    return expenses
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

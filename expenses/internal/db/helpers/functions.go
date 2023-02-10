package dbHelpers

// imports
import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

func NewExpense(ownerid string, name string, spent int) (Response, error) {
	id := uuid.New().String()
	_, err := DB.Exec("INSERT INTO expenses (owner_id, name, spent, id) VALUES ($1, $2, $3, $4)", ownerid, name, spent, id)
	response := Response{
		ID: id,
	}
	// return expenses
	return response, err
}
func DeleteExpense(id string) (error) {
	var err error
	if _, err := GetExpenseById(id); err != nil {
		return  err
	}
	_, err = DB.Exec(`DELETE FROM expenses WHERE id=$1`, id)
	return err
}
func UpdateExpense(id string, name string, spent int) (Response, error) {
    _, err := DB.Exec(`UPDATE expenses SET date_updated=now(), name=$1, spent=$2 WHERE id=$3`, name, spent, id)
	response := Response{
		ID: id,
	}
	return response, err
}

// apply changes to db (no breaking ones)
func Migrate() {
	fmt.Println("Migrating...")
	DB.MustExec(defaultSchema.create)
	ExecMultiple(DB, defaultSchema.alter)
	fmt.Println("Migrated!!")
}

// WARNING: THIS FUNCTION RESETS THE DATABASE
func ResetToSchema() {
	fmt.Println("Resetting...")
	ExecMultiple(DB, defaultSchema.drop)
	DB.MustExec(defaultSchema.create)
	ExecMultiple(DB, defaultSchema.alter)
	fmt.Println("Resetted!!")
}
func GetExpenseById(expenseId string) (Expense, error) {
	expense := Expense{}
	err := DB.Get(&expense, "SELECT * FROM expenses where id=$1", expenseId)
	if err == sql.ErrNoRows {
		return expense, ErrExpenseNotFound
	}
	return expense, err
}
func GetExpensesByOwnerId(ownerid string) ([]Expense, error) {
	expenses := []Expense{}
	err := DB.Select(&expenses, "SELECT * FROM expenses where owner_id=$1", ownerid)
	if len(expenses) <= 0 || err == sql.ErrNoRows {
		return expenses, ErrExpensesNotFound
	}
	return expenses, err
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

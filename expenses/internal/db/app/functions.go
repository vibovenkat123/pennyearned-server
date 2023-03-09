package db

// imports
import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"strings"
    helpers "main/expenses/internal/db/pkg"
)

func NewExpense(ownerid string, name string, spent int) (helpers.Response, error) {
	id := uuid.New().String()
	_, err := helpers.DB.Exec("INSERT INTO expenses (owner_id, name, spent, id) VALUES ($1, $2, $3, $4)", ownerid, name, spent, id)
	response := helpers.Response{
		ID: id,
	}
	// return expenses
	return response, err
}
func DeleteExpense(id string) error {
	var err error
	if _, err := GetExpenseById(id); err != nil {
		return err
	}
	_, err = helpers.DB.Exec(`DELETE FROM expenses WHERE id=$1`, id)
	return err
}
func UpdateExpense(id string, name string, spent int) (helpers.Response, error) {
	_, err := helpers.DB.Exec(`UPDATE expenses SET date_updated=now(), name=$1, spent=$2 WHERE id=$3`, name, spent, id)
	response := helpers.Response{
		ID: id,
	}
	return response, err
}

// apply changes to db (no breaking ones)
func Migrate() {
	fmt.Println("Migrating...")
	helpers.DB.MustExec(helpers.DefaultSchema.Create)
	ExecMultiple(helpers.DB, helpers.DefaultSchema.Alter)
	fmt.Println("Migrated!!")
}

// WARNING: THIS FUNCTION RESETS THE DATABASE
func ResetToSchema() {
	fmt.Println("Resetting...")
	ExecMultiple(helpers.DB, helpers.DefaultSchema.Drop)
	helpers.DB.MustExec(helpers.DefaultSchema.Create)
	ExecMultiple(helpers.DB, helpers.DefaultSchema.Alter)
	fmt.Println("Resetted!!")
}
func GetExpenseById(expenseId string) (helpers.Expense, error) {
	expense := helpers.Expense{}
	err := helpers.DB.Get(&expense, "SELECT * FROM expenses where id=$1", expenseId)
	if err == sql.ErrNoRows {
		return expense, helpers.ErrExpenseNotFound
	}
	return expense, err
}
func GetExpensesByOwnerId(ownerid string) ([]helpers.Expense, error) {
	expenses := []helpers.Expense{}
	err := helpers.DB.Select(&expenses, "SELECT * FROM expenses where owner_id=$1", ownerid)
	if len(expenses) <= 0 || err == sql.ErrNoRows {
		return expenses, helpers.ErrExpensesNotFound
	}
	return expenses, err
}
func ExecMultiple(e helpers.DatabaseType, query string) {
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

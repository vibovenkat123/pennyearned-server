package db

// imports
import (
	"database/sql"
	helpers "main/expenses/internal/db/pkg"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitializeLogger(logger *zap.Logger) {
	Logger = logger
}

func NewExpense(ownerid string, name string, spent int) (helpers.IDResponse, error) {
	id := uuid.New().String()
	_, err := helpers.DB.Exec("INSERT INTO expenses (owner_id, name, spent, id) VALUES ($1, $2, $3, $4)", ownerid, name, spent, id)
	response := helpers.IDResponse{
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
func UpdateExpense(id string, name string, spent int) (helpers.IDResponse, error) {
	_, err := helpers.DB.Exec(`UPDATE expenses SET date_updated=now(), name=$1, spent=$2 WHERE id=$3`, name, spent, id)
	response := helpers.IDResponse{
		ID: id,
	}
	return response, err
}

// apply changes to db (no breaking ones)
func Migrate() {
	Logger.Info("Migrating...")
	helpers.DB.MustExec(helpers.DefaultSchema.Create)
	ExecMultiple(helpers.DB, helpers.DefaultSchema.Alter)
	Logger.Info("Migrated!!")
}

// WARNING: THIS FUNCTION RESETS THE DATABASE
func ResetToSchema() {
	Logger.Info("Resetting...")
	ExecMultiple(helpers.DB, helpers.DefaultSchema.Drop)
	helpers.DB.MustExec(helpers.DefaultSchema.Create)
	ExecMultiple(helpers.DB, helpers.DefaultSchema.Alter)
	Logger.Info("Resetted!!")
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
			Logger.Error("Error executing statements",
				zap.Error(err),
			)
		}
	}
}

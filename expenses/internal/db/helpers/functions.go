package dbHelpers

// imports
import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"strings"
)

//func GetAllExpenses() ([]Expense, error) {
//	var expenses []Expense
//	err := DB.Select(&expenses, "SELECT * FROM expenses ORDER BY date_created ASC")
//	return expenses, err
//}

//func MostExpensiveExpense(ownerid string) (Expense, error) {
//	expenses, err := GetExpensesByOwnerId(ownerid)
//	if err != nil {
//		return Expense{}, err
//	}
//	// first array items expense and id
//	mostExpensive := expenses[0].Spent
//	id := expenses[0].ID
//	// start at index 1
//	for _, curr := range expenses[1:] {
//		if curr.Spent > mostExpensive {
//			mostExpensive = curr.Spent
//			id = curr.ID
//		}
//	}
//	// return the expense
//	return GetExpenseById(id)
//}

//func LeastExpensiveExpense(ownerid string) (Expense, error) {
//	expenses, err := GetExpensesByOwnerId(ownerid)
//	if err != nil {
//		return Expense{}, err
//	}
//	// first array items expense and id
//	leastExpensive := expenses[0].Spent
//	id := expenses[0].ID
//	// start at index 1
//	for _, curr := range expenses[1:] {
//		if curr.Spent < leastExpensive {
//			leastExpensive = curr.Spent
//			id = curr.ID
//		}
//	}
//	// return the expense
//	return GetExpenseById(id)
//}

//func ExpensesLowerThan(spent int, ownerid string) ([]Expense, error) {
//	var expensesLowerThan []Expense
//	expenses, err := GetExpensesByOwnerId(ownerid)
//	if err != nil {
//		return nil, err
//	}
//	for _, curr := range expenses {
//		if curr.Spent < spent {
//			expensesLowerThan = append(expensesLowerThan, curr)
//		}
//	}
//    if len(expenses) <= 0 {
//        return expensesLowerThan, ErrExpensesNotFound
//    }
//	return expensesLowerThan, nil
//}

//func ExpensesHigherThan(spent int, ownerid string) ([]Expense, error) {
//	var expensesHigherThan []Expense
//	expenses, err := GetExpensesByOwnerId(ownerid)
//	if err != nil {
//		return nil, err
//	}
//	for _, curr := range expenses {
//		if curr.Spent > spent {
//			expensesHigherThan = append(expensesHigherThan, curr)
//		}
//	}
//    if len(expenses) <= 0 {
//        return expensesHigherThan, ErrExpensesNotFound
//    }
//	return expensesHigherThan, nil
//}

// new expenses
func NewExpense(ownerid string, name string, spent int) (Response, error) {
	id := uuid.New().String()
	_, err := DB.Exec("INSERT INTO expenses (owner_id, name, spent, id) VALUES ($1, $2, $3, $4)", ownerid, name, spent, id)
	response := Response{
		ID: id,
	}
	// return expenses
	return response, err
}
func DeleteExpense(id string) (Response, error) {
	var err error
	if _, err := GetExpenseById(id); err != nil {
		return Response{ID: "Invalid Input"}, err
	}
	_, err = DB.Exec(`DELETE FROM expenses WHERE id=$1`, id)
	response := Response{
		ID: id,
	}
	return response, err
}
func UpdateExpense(id string, inputName string, inputSpent string) (Response, error) {
	original, err := GetExpenseById(id)
	name := inputName
	var spent int
	if len(name) <= 0 {
		name = original.Name
	}
	if len(inputSpent) <= 0 {
		spent = original.Spent
	} else {
		spent, err = strconv.Atoi(inputSpent)
		if err != nil {
			return Response{}, err
		}
	}
	// check if the inputs are nil
	// if they are, set them to the original
	_, err = DB.Exec(`UPDATE expenses SET date_updated=now(), name=$1, spent=$2 WHERE id=$3`, name, spent, id)
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

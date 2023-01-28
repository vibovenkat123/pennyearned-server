package dbHelpers

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"log"
	"strings"
)

func generateFromPassword(pwd string, p *params) (encodedHash string, err error) {
	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(pwd), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)
	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
func SignIn(email string, password string) (User, error) {
	user := User{}
	emptyUser := User{}
	Db.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	if user == emptyUser {
		return emptyUser, ErrEmailNotFound
	}
	matches, err := comparePasswordAndHash(password, user.Password)
	if err != nil {
		return emptyUser, err
	}
	if !matches {
		err = ErrPassNotMatch
		return emptyUser, err
	}
	user = User{
		ID:          user.ID,
		Name:        user.Name,
		Username:    user.Username,
		Password:    password,
		Email:       user.Email,
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
	}
	return user, nil
}
func SignUp(name string, username string, email string, password string) (User, error) {
	user := User{}
	id := uuid.New().String()
	encodedHash, err := generateFromPassword(password, p)
	if err != nil {
		return user, err
	}
	tx := Db.MustBegin()
	tx.MustExec("INSERT INTO users (name, username, email, password, id) VALUES ($1, $2, $3, $4, $5)", name, username, email, encodedHash, id)
	tx.Commit()
	Db.Get(&user, "SELECT * FROM users where email=$1", email)
	return user, nil
}
func decodeHash(encodedHash string) (p *params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
func comparePasswordAndHash(password string, encodedHash string) (matches bool, err error) {
	_, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}
	inputHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	if subtle.ConstantTimeCompare(hash, inputHash) == 1 {
		return true, nil
	}
	return false, nil
}
func GetAllExpenses() []Expense {
	var expenses []Expense
	Db.Select(&expenses, "SELECT * FROM expenses ORDER BY date_created ASC")
	return expenses
}
func MostExpensiveExpense(ownerid string) Expense {
	expensesData := GetExpensesByOwnerId(ownerid)
	mostExpensive := expensesData[0].Spent
	id := expensesData[0].ID
	for _, i := range expensesData[1:] {
		if i.Spent > mostExpensive {
			mostExpensive = i.Spent
			id = i.ID
		}
	}
	return GetExpenseById(id)
}
func LeastExpensiveExpense(ownerid string) Expense {
	expensesData := GetExpensesByOwnerId(ownerid)
	leastExpensive := expensesData[0].Spent
	id := expensesData[0].ID
	for _, i := range expensesData[1:] {
		if i.Spent < leastExpensive {
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
		if i.Spent < spent {
			expenses = append(expenses, i)
		}
	}
	return expenses
}
func ExpensesHigherThan(spent int, ownerid string) []Expense {
	var expenses []Expense
	expensesData := GetExpensesByOwnerId(ownerid)
	for _, i := range expensesData {
		if i.Spent > spent {
			expenses = append(expenses, i)
		}
	}
	return expenses
}
func NewExpense(ownerid string, name string, spent int) Expense {
	tx := Db.MustBegin()
	id := uuid.New().String()
	tx.MustExec("INSERT INTO expenses (owner_id, name, spent, id) VALUES ($1, $2, $3, $4)", ownerid, name, spent, id)
	tx.Commit()
	return GetExpenseById(id)
}
func DeleteExpense(id string) Expense {
	expense := GetExpenseById(id)
	Db.MustExec(`DELETE FROM expenses WHERE id=$1`, id)
	return expense
}
func UpdateExpense(expense Expense) Expense {
	Db.MustExec(`UPDATE expenses SET date_updated=now(), owner_id=$1, name=$2, spent=$3 WHERE id=$4`, expense.OwnerID, expense.Name, expense.Spent, expense.ID)
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

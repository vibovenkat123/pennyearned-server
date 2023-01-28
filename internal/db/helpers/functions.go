package dbHelpers

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	pq "github.com/lib/pq"
	"golang.org/x/crypto/argon2"
	"log"
	"net/mail"
	"strings"
)

func GenerateFromPassword(pwd string, p *params) (encodedHash string, err error) {
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
	err := Db.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	if err != nil && err != sql.ErrNoRows {
		return user, err
	}
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
		Password:    user.Password,
		Email:       user.Email,
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
	}
	return user, nil
}
func validEmailAddress(email string) (string, bool) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", false
	}
	return addr.Address, true
}
func validate(name string, username string, email string, password string) (string, error) {
	email, valid := validEmailAddress(email)
	if !valid {
		return "", ErrEmailInvalid
	}
	if len(name) <= 0 {
		return email, ErrNameTooShort
	}
	if len(name) >= 30 {
		return email, ErrNameTooLong
	}
	if len(username) <= 2 {
		return email, ErrUsernameTooShort
	}
	if len(username) >= 12 {
		return email, ErrUsernameTooLong
	}
	if len(password) < 8 {
		return email, ErrPasswordTooShort
	}
	if len(password) >= 20 {
		return email, ErrPasswordTooLong
	}
	return email, nil
}
func SignUp(name string, username string, email string, password string) (User, error) {
	user := User{}
	email, err := validate(name, username, email, password)
	if err != nil {
		return user, err
	}
	id := uuid.New().String()
	encodedHash, err := GenerateFromPassword(password, P)
	if err != nil {
		return user, err
	}
	_, err = Db.Exec("INSERT INTO users (name, username, email, password, id) VALUES ($1, $2, $3, $4, $5)", name, username, email, encodedHash, id)
	if err != nil {
		if err, _ := err.(*pq.Error); err.Code == "23505" {
			return user, ErrEmailAlreadyFound
		}
		return user, err
	}
	err = Db.Get(&user, "SELECT * FROM users where email=$1", email)
	return user, err
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
	inputHash := argon2.IDKey([]byte(password), salt, P.iterations, P.memory, P.parallelism, P.keyLength)

	if subtle.ConstantTimeCompare(hash, inputHash) == 1 {
		return true, nil
	}
	return false, nil
}
func GetAllExpenses() ([]Expense, error) {
	var expenses []Expense
	err := Db.Select(&expenses, "SELECT * FROM expenses ORDER BY date_created ASC")
	return expenses, err
}
func MostExpensiveExpense(ownerid string) (Expense, error) {
	expensesData, err := GetExpensesByOwnerId(ownerid)
	if err != nil {
		return Expense{}, err
	}
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
func LeastExpensiveExpense(ownerid string) (Expense, error) {
	expensesData, err := GetExpensesByOwnerId(ownerid)
	if err != nil {
		return Expense{}, err
	}
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
func ExpensesLowerThan(spent int, ownerid string) ([]Expense, error) {
	var expenses []Expense
	expensesData, err := GetExpensesByOwnerId(ownerid)
	if err != nil {
		return nil, err
	}
	for _, i := range expensesData {
		if i.Spent < spent {
			expenses = append(expenses, i)
		}
	}
	return expenses, nil
}
func ExpensesHigherThan(spent int, ownerid string) ([]Expense, error) {
	var expenses []Expense
	expensesData, err := GetExpensesByOwnerId(ownerid)
	if err != nil {
		return nil, err
	}
	for _, i := range expensesData {
		if i.Spent > spent {
			expenses = append(expenses, i)
		}
	}
	return expenses, nil
}
func NewExpense(ownerid string, name string, spent int) (Expense, error) {
	id := uuid.New().String()
	_, err := Db.Exec("INSERT INTO expenses (owner_id, name, spent, id) VALUES ($1, $2, $3, $4)", ownerid, name, spent, id)
	if err != nil {
		return Expense{}, err
	}
	return GetExpenseById(id)
}
func DeleteExpense(id string) (Expense, error) {
	expense, err := GetExpenseById(id)
	Db.MustExec(`DELETE FROM expenses WHERE id=$1`, id)
	return expense, err
}
func UpdateExpense(expense Expense) (Expense, error) {
	Db.MustExec(`UPDATE expenses SET date_updated=now(), owner_id=$1, name=$2, spent=$3 WHERE id=$4`, expense.OwnerID, expense.Name, expense.Spent, expense.ID)
	return GetExpenseById(expense.ID)
}
func GetUserById(id string) (User, error) {
	user := User{}
	err := Db.Get(&user, "SELECT * FROM users where id=$1", id)
	return user, err
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
func GetExpenseById(expenseId string) (Expense, error) {
	expense := Expense{}
	err := Db.Get(&expense, "SELECT * FROM expenses where id=$1", expenseId)
	return expense, err
}
func GetExpensesByOwnerId(ownerid string) ([]Expense, error) {
	expenses := []Expense{}
	err := Db.Select(&expenses, "SELECT * FROM expenses where owner_id=$1", ownerid)
	return expenses, err
}
func GetOwnerById(ownerid string) (User, error) {
	user := User{}
	err := Db.Select(&user, "SELECT * FROM users where _id=$1", ownerid)
	return user, err
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

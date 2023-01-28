package dbHelpers

// imports
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

// has a password and return it
func GenerateFromPassword(pwd string, p *params) (encodedHash string, err error) {
	// get the salt of the password
	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return "", err
	}
	// hash it
	hash := argon2.IDKey([]byte(pwd), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	// make the hash a string (base64)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	// return the encoded hash
	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)
	return encodedHash, nil
}

// generate a salt
func generateRandomBytes(n uint32) ([]byte, error) {
	// make a byte with the given salt length
	b := make([]byte, n)
	// generate the random bytes
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
func SignIn(email string, password string) (User, error) {
	// for returning
	user := User{}
	// empty struct for returning
	emptyUser := User{}
	// check if the email address is valid
	email, valid := validEmailAddress(email)
	if !valid {
		return user, ErrEmailInvalid
	}
	// get the user with the email address
	err := DB.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		// if there is no rows
		if err == sql.ErrNoRows {
			return user, ErrEmailNotFound
		}
		return user, err
	}
	// check if the password matches
	matches, err := comparePasswordAndHash(password, user.Password)
	if err != nil {
		return emptyUser, err
	}
	if !matches {
		err = ErrPassNotMatch
		return emptyUser, err
	}
	// create the user struct
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

// check if the email address is valid
func validEmailAddress(email string) (string, bool) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", false
	}
	return addr.Address, true
}

// validate all the fields
func validate(name string, username string, email string, password string) (string, error) {
	email, valid := validEmailAddress(email)
	if !valid {
		return "", ErrEmailInvalid
	}
	if len(name) <= minNameLength {
		return email, ErrNameTooShort
	}
	if len(name) >= maxNameLength {
		return email, ErrNameTooLong
	}
	if len(username) <= minUsernameLength {
		return email, ErrUsernameTooShort
	}
	if len(username) >= maxUsernameLength {
		return email, ErrUsernameTooLong
	}
	if len(password) <= minPasswordLength{
		return email, ErrPasswordTooShort
	}
	if len(password) >= maxPasswordLength {
		return email, ErrPasswordTooLong
	}
	return email, nil
}

// sign up
func SignUp(name string, username string, email string, password string) (User, error) {
	user := User{}
	// validate the email
	email, err := validate(name, username, email, password)
	if err != nil {
		return user, err
	}
	// generate random uuid
	id := uuid.New().String()
	// has the password
	encodedHash, err := GenerateFromPassword(password, P)
	if err != nil {
		return user, err
	}
	// insert a user
	_, err = DB.Exec("INSERT INTO users (name, username, email, password, id) VALUES ($1, $2, $3, $4, $5)", name, username, email, encodedHash, id)
	if err != nil {
		// if the error is "unique key violation" (email or username already found)
		if err, _ := err.(*pq.Error); err.Code == "23505" {
			return user, ErrAlreadyFound
		}
		return user, err
	}
	// get the user
	err = DB.Get(&user, "SELECT * FROM users where email=$1", email)
	return user, err
}

// decode the hash
func decodeHash(encodedHash string) (p *params, salt, hash []byte, err error) {
	// get the vals
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}
	// variable for version
	var version int
	// scan the values and put them in version
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	// if the versions dont match
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}
	// put the parameters in a struct
	p = &params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}
	// get the salt
	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	// get the salt length
	p.saltLength = uint32(len(salt))
	// get the hash
	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	// get keylength
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}

// check if password matches
func comparePasswordAndHash(password string, encodedHash string) (matches bool, err error) {
	// decode the hash
	_, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}
	// get the input hash
	inputHash := argon2.IDKey([]byte(password), salt, P.iterations, P.memory, P.parallelism, P.keyLength)
	// if they match
	if subtle.ConstantTimeCompare(hash, inputHash) == 1 {
		return true, nil
	}
	return false, nil
}
func GetAllExpenses() ([]Expense, error) {
	var expenses []Expense
	err := DB.Select(&expenses, "SELECT * FROM expenses ORDER BY date_created ASC")
	return expenses, err
}
func MostExpensiveExpense(ownerid string) (Expense, error) {
	expenses, err := GetExpensesByOwnerId(ownerid)
	if err != nil {
		return Expense{}, err
	}
	// first array items expense and id
	mostExpensive := expenses[0].Spent
	id := expenses[0].ID
	// start at index 1
	for _, curr := range expenses[1:] {
		if curr.Spent > mostExpensive {
			mostExpensive = curr.Spent
			id = curr.ID
		}
	}
	// return the expense
	return GetExpenseById(id)
}
func LeastExpensiveExpense(ownerid string) (Expense, error) {
	expenses, err := GetExpensesByOwnerId(ownerid)
	if err != nil {
		return Expense{}, err
	}
	// first array items expense and id
	leastExpensive := expenses[0].Spent
	id := expenses[0].ID
	// start at index 1
	for _, curr := range expenses[1:] {
		if curr.Spent < leastExpensive {
			leastExpensive = curr.Spent
			id = curr.ID
		}
	}
	// return the expense
	return GetExpenseById(id)
}

// get all expenses lower than
func ExpensesLowerThan(spent int, ownerid string) ([]Expense, error) {
	var expensesLowerThan []Expense
	expenses, err := GetExpensesByOwnerId(ownerid)
	if err != nil {
		return nil, err
	}
	for _, curr := range expenses {
		if curr.Spent < spent {
			expensesLowerThan = append(expensesLowerThan, curr)
		}
	}
	return expensesLowerThan, nil
}

// return all expenses higher than
func ExpensesHigherThan(spent int, ownerid string) ([]Expense, error) {
	var expensesHigherThan []Expense
	expenses, err := GetExpensesByOwnerId(ownerid)
	if err != nil {
		return nil, err
	}
	for _, curr := range expenses {
		if curr.Spent > spent {
			expensesHigherThan = append(expenses, curr)
		}
	}
	return expensesHigherThan, nil
}

// new expenses
func NewExpense(ownerid string, name string, spent int) (Expense, error) {
	id := uuid.New().String()
	_, err := DB.Exec("INSERT INTO expenses (owner_id, name, spent, id) VALUES ($1, $2, $3, $4)", ownerid, name, spent, id)
	if err != nil {
		// return empty expense
		return Expense{}, err
	}
	// return expenses
	return GetExpenseById(id)
}
func DeleteExpense(id string) (Expense, error) {
	expense, err := GetExpenseById(id)
	if err != nil {
		return expense, err
	}
	_, err = DB.Exec(`DELETE FROM expenses WHERE id=$1`, id)
	return expense, err
}
func UpdateExpense(expense Expense) (Expense, error) {
	_, err := DB.Exec(`UPDATE expenses SET date_updated=now(), owner_id=$1, name=$2, spent=$3 WHERE id=$4`, expense.OwnerID, expense.Name, expense.Spent, expense.ID)
	if err != nil {
		// return empty expense
		return Expense{}, err
	}
	return GetExpenseById(expense.ID)
}
func GetUserById(id string) (User, error) {
	user := User{}
	err := DB.Get(&user, "SELECT * FROM users where id=$1", id)
	return user, err
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
	return expense, err
}
func GetExpensesByOwnerId(ownerid string) ([]Expense, error) {
	expenses := []Expense{}
	err := DB.Select(&expenses, "SELECT * FROM expenses where owner_id=$1", ownerid)
	return expenses, err
}
func GetOwnerById(ownerid string) (User, error) {
	user := User{}
	err := DB.Select(&user, "SELECT * FROM users where _id=$1", ownerid)
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

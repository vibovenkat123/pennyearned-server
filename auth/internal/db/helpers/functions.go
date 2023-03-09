package dbHelpers

// imports
import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"log"
	"net/http"
	"strings"
	"time"
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
func GetByCookie(cookieID string, ctx context.Context) (string, error) {
	val, err := RDB.Get(ctx, fmt.Sprintf("session:%v", cookieID)).Result()
	return val, err
}
func SignOut(cookieID string, ctx context.Context) (cookie *http.Cookie, err error) {
	// remove the cookie
	cookie = &http.Cookie{
		Name:    "user",
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
	}
	// check if the cookie exists in database
	_, err = RDB.Get(ctx, fmt.Sprintf("session:%v", cookieID)).Result()
	if err != nil {
		return cookie, err
	}
	// delete the entry
	err = RDB.Del(ctx, fmt.Sprintf("session:%v", cookieID)).Err()
	if err != nil {
		return cookie, err
	}
	return cookie, err
}
func SignIn(email string, password string, ctx context.Context) (*http.Cookie, error) {
	user := User{}
	err := DB.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	// if the database returns no rows
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	// if the password and hash match
	matches, err := comparePasswordAndHash(password, user.Password)
	if err != nil || !matches {
		if !matches {
			err = ErrPassNotMatch
		}
		return nil, err
	}
	// create the cookie
	cookie, cookieID, expireTime := CreateUserCookie()
	err = RDB.Set(ctx, fmt.Sprintf("session:%v", cookieID), user.ID, expireTime).Err()
	return cookie, err
}
func CreateUserCookie() (*http.Cookie, string, time.Duration) {
	expireTime := 86400 * time.Second
	expires := time.Now().Add(expireTime)
	cookieID := uuid.New().String()
	cookie := &http.Cookie{
		Name:    "user",
		Value:   cookieID,
		Expires: expires,
		Path:    "/",
	}
	return cookie, cookieID, expireTime
}

func SignUp(name string, username string, email string, password string, ctx context.Context) (*http.Cookie, error) {
	id := uuid.New().String()
	encodedHash, err := GenerateFromPassword(password, P)
	if err != nil {
		return nil, err
	}
	_, err = DB.Exec("INSERT INTO users (name, username, email, password, id) VALUES ($1, $2, $3, $4, $5)", name, username, email, encodedHash, id)
	if err != nil {
		return nil, err
	}
	cookie, cookieID, expireTime := CreateUserCookie()
	err = RDB.Set(ctx, fmt.Sprintf("session:%v", cookieID), id, expireTime).Err()
	return cookie, err
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

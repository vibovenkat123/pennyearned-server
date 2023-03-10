package db

// imports
import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
	helpers "main/auth/internal/db/pkg"
	"net/http"
	"strings"
	"time"
)

var log *zap.Logger

func InitializeLogger(logger *zap.Logger) {
	log = logger
}
func GenerateFromPassword(pwd string, p *helpers.Params) (encodedHash string, err error) {
	salt, err := generateRandomBytes(p.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(pwd), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)
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
	val, err := helpers.RDB.Get(ctx, fmt.Sprintf("session:%v", cookieID)).Result()
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
	_, err = helpers.RDB.Get(ctx, fmt.Sprintf("session:%v", cookieID)).Result()
	if err != nil {
		return cookie, err
	}
	// delete the entry
	err = helpers.RDB.Del(ctx, fmt.Sprintf("session:%v", cookieID)).Err()
	if err != nil {
		return cookie, err
	}
	return cookie, err
}
func SignIn(email string, password string, ctx context.Context) (*http.Cookie, error) {
	user := helpers.User{}
	err := helpers.DB.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	// if the database returns no rows
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	// if the password and hash match
	matches, err := comparePasswordAndHash(password, user.Password)
	if err != nil || !matches {
		if !matches {
			err = helpers.ErrPassNotMatch
		}
		return nil, err
	}
	// create the cookie
	cookie, cookieID, expireTime := CreateUserCookie()
	err = helpers.RDB.Set(ctx, fmt.Sprintf("session:%v", cookieID), user.ID, expireTime).Err()
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
	encodedHash, err := GenerateFromPassword(password, helpers.P)
	if err != nil {
		return nil, err
	}
	_, err = helpers.DB.Exec("INSERT INTO users (name, username, email, password, id) VALUES ($1, $2, $3, $4, $5)", name, username, email, encodedHash, id)
	if err != nil {
		return nil, err
	}
	cookie, cookieID, expireTime := CreateUserCookie()
	err = helpers.RDB.Set(ctx, fmt.Sprintf("session:%v", cookieID), id, expireTime).Err()
	return cookie, err
}
func decodeHash(encodedHash string) (p *helpers.Params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, helpers.ErrInvalidHash
	}
	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, helpers.ErrIncompatibleVersion
	}

	p = &helpers.Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}
func comparePasswordAndHash(password string, encodedHash string) (matches bool, err error) {
	_, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}
	inputHash := argon2.IDKey([]byte(password), salt, helpers.P.Iterations, helpers.P.Memory, helpers.P.Parallelism, helpers.P.KeyLength)

	if subtle.ConstantTimeCompare(hash, inputHash) == 1 {
		return true, nil
	}
	return false, nil
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
func ExecMultiple(e helpers.DatabaseType, query string) {
	statements := strings.Split(query, "\n")
	if len(strings.Trim(statements[len(statements)-1], " \n\t\r")) == 0 {
		statements = statements[:len(statements)-1]
	}
	for _, s := range statements {
		_, err := e.Exec(s)
		if err != nil {
			log.Error("Error executing statements",
				zap.Error(err),
			)
		}
	}
}

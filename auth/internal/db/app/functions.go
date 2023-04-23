package db

// imports
import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	helpers "main/auth/internal/db/pkg"
	"net/smtp"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
	"os"
)

var log *zap.Logger
var tmpt *template.Template
var fromEmail string
var emailPassword string
var smtpHost string
var smtpPort string
var auth smtp.Auth
var templateFile string

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

type EmailData struct {
	Code string
}

func init() {
	if os.Getenv("GO_ENV") == "local" {
		fromEmail = os.Getenv("FROM_EMAIL")
		emailPassword = os.Getenv("FROM_PASSWORD")
		smtpHost = os.Getenv("SMTP_HOST")
		smtpPort = os.Getenv("SMTP_PORT")
		templateFile = os.Getenv("TEMPLATE_PATH")
	}
	tmpt = template.Must(template.ParseFiles(templateFile))
	auth = smtp.PlainAuth("", fromEmail, emailPassword, smtpHost)
}
func InitializeLogger(logger *zap.Logger) {
	log = logger
}
func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func SendEmail(to []string, ctx context.Context) error {
	for _, email := range to {
		code := EncodeToString(6)
		log.Info("Generated verification code")
		log.Debug(code)

		// create a new body buffer for each email
		body := bytes.NewBuffer(nil)
		mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		body.Write([]byte(fmt.Sprintf("Subject: %v is your verification code\n%s\n\n", code, mimeHeaders)))
		data := EmailData{
			Code: code,
		}
		log.Info(code)
		err := helpers.RDB.Set(ctx, fmt.Sprintf("code:%v", code), email, time.Second*1800).Err()
		if err != nil {
			log.Error("Error setting code in the redis",
				zap.Error(err),
			)
			return err
		}
		err = tmpt.Execute(body, data)
		if err != nil {
			return err
		}
		log.Info("Sending mail")
		err = smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{email}, body.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
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
func GetByAccess(accessToken string, ctx context.Context) (helpers.User, error) {
	user := helpers.User{}
	log.Info("getting session (cookie)")
	val, err := helpers.RDB.Get(ctx, fmt.Sprintf("session:%v", accessToken)).Result()
	if err != nil {
		log.Error(err.Error())
		return user, err
	}
	err = helpers.DB.Get(&user, "SELECT * FROM users WHERE id=?", val)
	if err != nil {
		log.Error(err.Error())
	}
	return user, err
}
func SignOut(accessToken string, ctx context.Context) (err error) {
	// remove the accessToken
	// check if the accessToken exists in database
	log.Info("Checking if the access token exists")
	_, err = helpers.RDB.Get(ctx, fmt.Sprintf("session:%v", accessToken)).Result()
	if err != nil {
		return err
	}
	// delete the entry
	err = helpers.RDB.Del(ctx, fmt.Sprintf("session:%v", accessToken)).Err()
	if err != nil {
		return err
	}
	return err
}
func SignIn(username string, password string, ctx context.Context) (accessToken *string, err error) {
	user := helpers.User{}
	log.Info("Getting the user",
		zap.String("username", username),
	)
	err = helpers.DB.Get(&user, "SELECT * FROM users WHERE username=?", username)
	// if the database returns no rows
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	// if the password and hash match
	matches, err := comparePasswordAndHash(password, user.Password)
	if err != nil || !matches {
		if !matches {
			log.Error("The password is invalid")
			err = helpers.ErrPassNotMatch
		}
		return nil, err
	}
	// create the accessToken
	accessToken, expireTime := CreateAccessToken()
	err = helpers.RDB.Set(ctx, fmt.Sprintf("session:%v", *accessToken), user.ID, expireTime).Err()
	return
}

func CreateAccessToken() (*string, time.Duration) {
	token := uuid.New().String()
	accessToken := &token
	return accessToken, time.Second * 86400
}

func SignUp(username string, password string, code string, ctx context.Context) (accessToken *string, err error) {
	id := uuid.New().String()
	encodedHash, err := GenerateFromPassword(password, helpers.P)
	if err != nil {
		return nil, err
	}
	log.Info("Getting the access token")
	email, err := helpers.RDB.Get(ctx, fmt.Sprintf("code:%v", code)).Result()
	log.Debug("Got the email",
		zap.String("email", email),
	)
	if email == "" {
		return nil, helpers.ErrInvalidCode
	}
	if err != nil {
		return nil, err
	}

	log.Info("Creating new user")
	_, err = helpers.DB.Exec("INSERT INTO users (username, email, password, id) VALUES (?, ?, ?, ?)", username, email, encodedHash, id)
	if err != nil {
		return nil, err
	}
	accessToken, expireTime := CreateAccessToken()
	log.Info("Creating new session")
	err = helpers.RDB.Set(ctx, fmt.Sprintf("session:%v", *accessToken), id, expireTime).Err()
	if err != nil {
		return accessToken, err
	}
	log.Info("Deleting existing code")
	err = helpers.RDB.Del(ctx, fmt.Sprintf("code:%v", code)).Err()
	return accessToken, err
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
	fmt.Println(inputHash)

	if subtle.ConstantTimeCompare(hash, inputHash) == 1 {
		return true, nil
	}
	return false, nil
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

package userFunctions;
import (
    "net/mail"
)
func Validate(name string, username string, email string, password string) (good bool) {
    return NameCheck(name) && UsernameCheck(username) && EmailCheck(email) && PasswordCheck(password) 
}
func NameCheck(name string) bool{
    length := len(name)
    return length > 0 && length < 30
}
func UsernameCheck(username string) bool{
    length := len(username)
    return length > 2 && length < 12
}
func EmailCheck(email string) bool{
    _, err := mail.ParseAddress(email)
    return err == nil
}
func PasswordCheck(password string) bool{
    length := len(password)
    return length >= 8 && length <= 20
}

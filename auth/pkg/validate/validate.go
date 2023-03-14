package validate

import "strconv"

func All(name string, username string, code string, password string) (good bool) {
	return Name(name) && Username(username) && Password(password) && Code(code)
}

func Code(code string) bool {
	_, err := strconv.Atoi(code)
	return err == nil && len(code) == 6
}

func Name(name string) bool {
	length := len(name)
	return length > 0 && length < 30
}
func Username(username string) bool {
	length := len(username)
	return length > 2 && length < 12
}
func Password(password string) bool {
	length := len(password)
	return length >= 8 && length <= 20
}

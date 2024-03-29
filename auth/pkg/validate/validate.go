package validate

import "strconv"

func All(username string, code string, password string) (good bool) {
	return Username(username) && Password(password) && Code(code)
}

func Code(code string) bool {
	_, err := strconv.Atoi(code)
	return err == nil && len(code) == 6
}

func Username(username string) bool {
	length := len(username)
	return length > 2 && length < 12
}
func Password(password string) bool {
	length := len(password)
	return length >= 8 && length <= 20
}

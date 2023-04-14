package validate

import (
	"testing"
)

func TestValidate(t *testing.T) {
	t.Run("Invalid", TestInvalidCases)
	t.Run("Valid", TestValidCases)
}

func TestValidCases(t *testing.T) {
	name := "bob"
	username := "bob123"
	code := "123456"
	password := "12345678"
	good := All(username, code, password)
	if !good {
		t.Fatalf(`All(%v, %v, %v, %v) returned %t with good cases, wanted it to return %t`, name, username, code, password, good, !good)
	}
}

func TestInvalidCases(t *testing.T) {
	t.Run("Invalid Password", TestInvalidPassword)
	t.Run("Invalid Code", TestInvalidCode)
	t.Run("Invalid Username", TestInvalidUsername)
}


func TestInvalidUsername(t *testing.T) {
	t.Run("Username too short", TestUsernameTooShort)
	t.Run("Username too long", TestUsernameTooLong)
}

func TestInvalidPassword(t *testing.T) {
	t.Run("Password too short", TestPasswordTooShort)
	t.Run("Password too long", TestPasswordTooLong)
}

func TestInvalidCode(t *testing.T) {
	email := "12345"
	good := Code(email)
	if good {
		t.Fatalf(`Code(%v) with a wrong code case returned %t, wanted return %t`, email, good, !good)
	}
}

func TestPasswordTooLong(t *testing.T) {
	password := "aaaaaaaaaaaaaaaaaaaaaaaaaa"
	good := Password(password)
	if good {
		t.Fatalf(`Password(%v) with a wrong password case returned %t, wanted return %t`, password, good, !good)
	}
}

func TestPasswordTooShort(t *testing.T) {
	password := "a"
	good := Password(password)
	if good {
		t.Fatalf(`Password(%v) with a wrong password case returned %t, wanted return %t`, password, good, !good)
	}
}

func TestUsernameTooLong(t *testing.T) {
	username := "aaaaaaaaaaaaaaaa"
	good := Username(username)
	if good {
		t.Fatalf(`Username(%v) with a wrong username case returned %t, wanted return %t`, username, good, !good)
	}
}

func TestUsernameTooShort(t *testing.T) {
	username := "aa"
	good := Username(username)
	if good {
		t.Fatalf(`Username(%v) with a wrong username case returned %t, wanted return %t`, username, good, !good)
	}
}


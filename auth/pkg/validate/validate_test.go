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
	email := "bob@gmail.com"
	password := "12345678"
	good := Validate(name, username, email, password)
	if !good {
		t.Fatalf(`Validate(%v, %v, %v, %v) returned %t with good cases, wanted it to return %t`, name, username, email, password, !good, good)
	}
}

func TestInvalidCases(t *testing.T) {
	t.Run("Invalid Name", TestInvalidName)
	t.Run("Invalid Password", TestInvalidPassword)
	t.Run("Invalid Email", TestInvalidEmail)
	t.Run("Invalid Username", TestInvalidUsername)
}

func TestInvalidName(t *testing.T) {
	t.Run("Name too short", TestNameTooShort)
	t.Run("Name too long", TestNameTooLong)
}

func TestInvalidUsername(t *testing.T) {
	t.Run("Username too short", TestUsernameTooShort)
	t.Run("Username too long", TestUsernameTooLong)
}

func TestInvalidPassword(t *testing.T) {
	t.Run("Password too short", TestPasswordTooShort)
	t.Run("Password too long", TestPasswordTooLong)
}

func TestInvalidEmail(t *testing.T) {
	email := "joe"
	good := EmailCheck(email)
	if good {
		t.Fatalf(`EmailCheck(%v) with a wrong email case returned %t, wanted return %t`, email, good, good)
	}
}

func TestPasswordTooLong(t *testing.T) {
	password := "aaaaaaaaaaaaaaaaaaaaaaaaaa"
	good := PasswordCheck(password)
	if good {
		t.Fatalf(`PasswordCheck(%v) with a wrong password case returned %t, wanted return %t`, password, good, good)
	}
}

func TestPasswordTooShort(t *testing.T) {
	password := "a"
	good := PasswordCheck(password)
	if good {
		t.Fatalf(`PasswordCheck(%v) with a wrong password case returned %t, wanted return %t`, password, good, good)
	}
}

func TestUsernameTooLong(t *testing.T) {
	username := "aaaaaaaaaaaaaaaa"
	good := UsernameCheck(username)
	if good {
		t.Fatalf(`UsernameCheck(%v) with a wrong username case returned %t, wanted return %t`, username, good, good)
	}
}

func TestUsernameTooShort(t *testing.T) {
	username := "aa"
	good := UsernameCheck(username)
	if good {
		t.Fatalf(`UsernameCheck(%v) with a wrong username case returned %t, wanted return %t`, username, good, good)
	}
}

func TestNameTooLong(t *testing.T) {
	name := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	good := NameCheck(name)
	if good {
		t.Fatalf(`NameCheck(%v) with a wrong name case returned %t, wanted return %t`, name, good, good)
	}
}

func TestNameTooShort(t *testing.T) {
	name := ""
	good := NameCheck(name)
	if good {
		t.Fatalf(`NameCheck(%v) with a wrong name case returned %t, wanted return %t`, name, good, good)
	}
}

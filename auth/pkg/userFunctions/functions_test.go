package userFunctions

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
		t.Fatalf(`Validate returned %t, wanted return %t`, good,  true)
	}
}

func TestInvalidCases(t *testing.T) {
	name := "bob"
    username := "bob123"
    email := "bob@gmail.com"
    password := "1234567"
	good := Validate(name, username, email, password)
	if good {
		t.Fatalf(`Validate returned %t, wanted return %t`, good,  false)
	}
}

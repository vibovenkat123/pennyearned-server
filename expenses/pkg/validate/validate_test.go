package validate

import (
	"testing"
)

func TestValidate(t *testing.T) {
	t.Run("Invalid", TestInvalidCases)
	t.Run("Valid", TestValidCases)
}

func TestValidCases(t *testing.T) {
	id := "025823d5-0f92-47f9-a182-ff1e79e987d6"
	name := "bob"
	spent := 123
	good := All(id, name, spent)
	if !good {
		t.Fatalf(`All(%v, %v, %v) returned %t with good cases, wanted it to return %v"`, id, name, spent, !good, good)
	}
}

func TestInvalidCases(t *testing.T) {
	t.Run("Invalid Spent", TestInvalidSpent)
	t.Run("Invalid ID", TestInvalidID)
	t.Run("Invalid Name", TestInvalidName)
}

func TestInvalidSpent(t *testing.T) {
	spent := -1
	good := Spent(spent)
	if good {
		t.Fatalf(`Spent(%v) with a wrong spent case returned %t, wanted it to return %v"`, spent, good, !good)
	}
}

func TestInvalidID(t *testing.T) {
	id := "123"
	good := ID(id)
	if good {
		t.Fatalf(`ID(%v) with a wrong id case returned %t, wanted it to return %v"`, id, good, !good)
	}
}

func TestInvalidName(t *testing.T) {
	t.Run("Name too short", TestNameTooShort)
	t.Run("Name too long", TestNameTooLong)
}
func TestNameTooLong(t *testing.T) {
	name := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	good := Name(name)
	if good {
		t.Fatalf(`Name(%v) with a wrong name case returned %t, wanted it to return %v"`, name, good, !good)
	}
}
func TestNameTooShort(t *testing.T) {
	name := ""
	good := Name(name)
	if good {
		t.Fatalf(`Name(%v) with a wrong name case returned %t, wanted it to return %v"`, name, good, !good)
	}
}

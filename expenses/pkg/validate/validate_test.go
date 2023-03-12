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
	good := Validate(id, name, spent)
	if !good {
		t.Fatalf(`Validate(%v, %v, %v) returned %t with good cases, wanted it to return %v"`, id, name, spent, !good, good)
	}
}

func TestInvalidCases(t *testing.T) {
	t.Run("Invalid Spent", TestInvalidSpent)
	t.Run("Invalid ID", TestInvalidID)
	t.Run("Invalid Name", TestInvalidName)
}

func TestInvalidSpent(t *testing.T) {
	spent := -1
	good := ValidateSpent(spent)
	if good {
		t.Fatalf(`ValidateSpent(%v) with a wrong spent case returned %t, wanted it to return %v"`, spent, good, !good)
	}
}

func TestInvalidID(t *testing.T) {
	id := "123"
	good := ValidateID(id)
	if good {
		t.Fatalf(`ValidateID(%v) with a wrong id case returned %t, wanted it to return %v"`, id, good, !good)
	}
}

func TestInvalidName(t *testing.T) {
	t.Run("Name too short", TestNameTooShort)
	t.Run("Name too long", TestNameTooLong)
}
func TestNameTooLong(t *testing.T) {
	name := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	good := ValidateName(name)
	if good {
		t.Fatalf(`ValidateName(%v) with a wrong name case returned %t, wanted it to return %v"`, name, good, !good)
	}
}
func TestNameTooShort(t *testing.T) {
	name := ""
	good := ValidateName(name)
	if good {
		t.Fatalf(`ValidateName(%v) with a wrong name case returned %t, wanted it to return %v"`, name, good, !good)
	}
}

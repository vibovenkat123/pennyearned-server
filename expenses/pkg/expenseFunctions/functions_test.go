package expenseFunctions

import (
	"testing"
)

func TestValidate(t *testing.T) {
	t.Run("Invalid", TestInvalidCases)
	t.Run("Valid", TestValidCases)
	t.Run("Invalid Spent", TestInvalidSpentCase)
	t.Run("Invalid ID", TestInvalidIDCase)
	t.Run("Invalid Name", TestInvalidNameCase)
}

func TestValidCases(t *testing.T) {
	id := "025823d5-0f92-47f9-a182-ff1e79e987d6"
	name := "bob"
	spent := 123
	good := Validate(id, name, spent)
	if !good {
		t.Fatalf(`Validate() = %t, want  result to be %v"`, good, true)
	}
}

func TestInvalidCases(t *testing.T) {
	wrongName := ""
	wrongSpent := -1
	wrongId := "123"
	good := Validate(wrongId, wrongName, wrongSpent)
	if good {
		t.Fatalf(`Validate() with wrong cases  = %t, want  result to be %v"`, good, false)
	}
}

func TestInvalidSpentCase(t *testing.T) {
	wrongSpent := -1
	good := ValidateSpent(wrongSpent)
	if good {
		t.Fatalf(`ValidateSpent() with wrong spent case returned %t, want  result to be %v"`, good, false)
	}
}

func TestInvalidIDCase(t *testing.T) {
	wrongID := "123"
	good := ValidateID(wrongID)
	if good {
		t.Fatalf(`ValidateID() with wrong id case returned %t, want  result to be %v"`, good, false)
	}
}

func TestInvalidNameCase(t *testing.T) {
	wrongName := ""
	good := ValidateName(wrongName)
	if good {
		t.Fatalf(`ValidateName() with wrong name case returned %t, want  result to be %v"`, good, false)
	}
}

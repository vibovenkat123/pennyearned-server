package expenseFunctions

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
	idIsGood, nameIsGood, spentIsGood := Validate(&id, &name, &spent)
	if !idIsGood {
		t.Fatalf(`Validate(nil, *%v, nil) = %t, want  result to be %v"`, name, nameIsGood, true)
	}
	if !nameIsGood {
		t.Fatalf(`Validate(nil, *%v, nil) = %t, want  result to be %v"`, id, idIsGood, true)
	}
	if !spentIsGood {
		t.Fatalf(`Validate(nil, *%v, nil) = %t, want  result to be %v"`, spent, spentIsGood, true)
	}
}

func TestInvalidCases(t *testing.T) {
	wrongName := ""
	wrongSpent := 1
	wrongId := "123"
	idIsGood, nameIsGood, spentIsGood := Validate(&wrongId, &wrongName, &wrongSpent)
	if idIsGood {
		t.Fatalf(`Validate(nil, *%v, nil) = %t, want  result to be %v"`, wrongId, idIsGood, false)
	}
	if nameIsGood {
		t.Fatalf(`Validate(nil, *%v, nil) = %t, want  result to be %v"`, wrongName, nameIsGood, false)
	}
	if spentIsGood {
		t.Fatalf(`Validate(nil, *%v, nil) = %t, want  result to be %v"`, wrongSpent, spentIsGood, false)
	}
}

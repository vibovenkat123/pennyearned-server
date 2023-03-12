package validate

func Validate(id string, name string, spent int) bool {
	return (ValidateID(id)) && (ValidateName(name)) && (ValidateSpent(spent))
}
func ValidateID(id string) bool {
	return len(id) == 36
}

func ValidateName(name string) bool {
	return len(name) > 1 && len(name) < 30
}

func ValidateSpent(spent int) bool {
	return spent > 0
}

package validate

func All(id string, name string, spent int) bool {
	return (ID(id)) && (Name(name)) && (Spent(spent))
}
func ID(id string) bool {
	return len(id) == 36
}

func Name(name string) bool {
	return len(name) > 1 && len(name) < 30
}

func Spent(spent int) bool {
	return spent >= 0
}

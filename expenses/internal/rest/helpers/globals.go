package apiHelpers;
var (
    Port = 3000
    GetExpenseByIDUrl = "/api/expense/get/"
    GetExpensesByOwnerIDUrl = "/api/expenses/get/"
)
type Error struct {
    Message string
    Code int
}

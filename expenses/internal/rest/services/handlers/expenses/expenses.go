package expenses;
import (
    "fmt"
    "net/http"
    "strings"
    "main/expenses/internal/db/helpers"
    "encoding/json"
    helpers "main/expenses/internal/rest/helpers"

)

func GetByID(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    encoder := json.NewEncoder(w)
    id := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%v", helpers.GetExpenseByIDUrl))
    expenses, err := dbHelpers.GetExpenseById(id) 
    if err != nil {
        errRes := helpers.Error {
            Message: err.Error(),
            Code: 200,
        }
        encoder.Encode(errRes)
    }else {
        encoder.Encode(expenses)
    }
}
func GetByOwnerID(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    encoder := json.NewEncoder(w)
    id := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%v", helpers.GetExpensesByOwnerIDUrl))
    expenses, err := dbHelpers.GetExpensesByOwnerId(id) 
    if err != nil {
        errRes := helpers.Error {
            Message: err.Error(),
            Code: 200,
        }
        encoder.Encode(errRes)
    }else {
        encoder.Encode(expenses)
    }
}

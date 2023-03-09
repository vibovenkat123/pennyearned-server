package expenses

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"main/expenses/internal/db/app"
	"main/expenses/pkg/expenseFunctions"
	"net/http"
	"strconv"
)

func ErrInvalidFormat(w http.ResponseWriter) {
	http.Error(w, http.StatusText(400), 400)
	w.WriteHeader(400)
}
func ErrNotFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(404), 404)
	w.WriteHeader(404)
}
func GetByID(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	id := chi.URLParam(r, "id")
	if i, _, _ := expenseFunctions.Validate(&id, nil, nil); !i {
		ErrInvalidFormat(w)
		return
	}
	expenses, err := db.GetExpenseById(id)
	if err != nil {
		ErrNotFound(w)
		return
	}
	encoder.Encode(expenses)
	w.WriteHeader(200)
}
func GetByOwnerID(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	ownerid := chi.URLParam(r, "id")
	if i, _, _ := expenseFunctions.Validate(&ownerid, nil, nil); !i {
		ErrInvalidFormat(w)
		return
	}
	expenses, err := db.GetExpensesByOwnerId(ownerid)
	if err != nil {
		ErrNotFound(w)
		return
	}
	encoder.Encode(expenses)
	w.WriteHeader(200)

}
func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	id := chi.URLParam(r, "id")
	name := r.URL.Query().Get("name")
	inputSpent := r.URL.Query().Get("spent")
	original, err := db.GetExpenseById(id)
	var spent int
	if len(name) <= 0 {
		name = original.Name
	}
	if len(inputSpent) <= 0 {
		spent = original.Spent
	} else {
		spent, err = strconv.Atoi(inputSpent)
		if err != nil {
			ErrInvalidFormat(w)
			return
		}
	}
	if i, n, s := expenseFunctions.Validate(&id, &name, &spent); !i || !n || !s {
		ErrInvalidFormat(w)
		return
	}
	response, err := db.UpdateExpense(id, name, spent)
	if err != nil {
		ErrNotFound(w)
		return
	}
	encoder.Encode(response)
	w.WriteHeader(200)
}
func NewExpense(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	ownerid := r.URL.Query().Get("ownerid")
	name := r.URL.Query().Get("name")
	spent, err := strconv.Atoi(r.URL.Query().Get("spent"))
	if err != nil {
		ErrInvalidFormat(w)
		return
	}
	i, n, s := expenseFunctions.Validate(&ownerid, &name, &spent)
	if !i || !n || !s {
		ErrInvalidFormat(w)
		return
	}
	response, err := db.NewExpense(ownerid, name, spent)
	if err != nil {
		ErrNotFound(w)
		return
	}
	encoder.Encode(response)
	w.WriteHeader(201)
}
func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if i, _, _ := expenseFunctions.Validate(&id, nil, nil); !i {
		ErrInvalidFormat(w)
		return
	}
	err := db.DeleteExpense(id)
	if err != nil {
		ErrNotFound(w)
		return
	}
	w.WriteHeader(204)
}

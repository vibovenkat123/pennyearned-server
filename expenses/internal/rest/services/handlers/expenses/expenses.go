package expenses

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"main/expenses/internal/db/helpers"
	"net/http"
	"strconv"
)

func GetByID(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	id := chi.URLParam(r, "id")
	if len(id) != 36 {
		http.Error(w, http.StatusText(400), 400)
		w.WriteHeader(400)
		return
	}
	expenses, err := dbHelpers.GetExpenseById(id)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		w.WriteHeader(404)
		return
	}
	encoder.Encode(expenses)
	w.WriteHeader(200)
}
func GetByOwnerID(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	ownerid := chi.URLParam(r, "id")
	if len(ownerid) != 36 {
		http.Error(w, http.StatusText(400), 400)
		w.WriteHeader(400)
		return
	}
	expenses, err := dbHelpers.GetExpensesByOwnerId(ownerid)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		w.WriteHeader(404)
		return
	}
	encoder.Encode(expenses)
	w.WriteHeader(200)

}
func validateExpenseInput(name string, ownerid string) bool {
	if len(name) <= 0 || len(ownerid) != 36 {
		return true
	}
	return false
}
func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	id := chi.URLParam(r, "id")
	name := r.URL.Query().Get("name")
	spent := r.URL.Query().Get("spent")
	if len(id) != 36 {
		http.Error(w, http.StatusText(400), 400)
		w.WriteHeader(400)
		return
	}
	response, err := dbHelpers.UpdateExpense(id, name, spent)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		w.WriteHeader(404)
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
	if err != nil || validateExpenseInput(name, ownerid) {
		http.Error(w, http.StatusText(400), 400)
		w.WriteHeader(400)
		return
	}
	response, err := dbHelpers.NewExpense(ownerid, name, spent)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		w.WriteHeader(404)
		return
	}
	encoder.Encode(response)
	w.WriteHeader(200)
}
func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	id := chi.URLParam(r, "id")
	if len(id) != 36 {
		http.Error(w, http.StatusText(400), 400)
		w.WriteHeader(400)
		return
	}
	response, err := dbHelpers.DeleteExpense(id)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		w.WriteHeader(404)
		return
	}
	encoder.Encode(response)
	w.WriteHeader(200)
}

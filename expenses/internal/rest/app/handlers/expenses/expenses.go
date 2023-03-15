package expenses

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"main/expenses/internal/db/app"
	globals "main/expenses/internal/rest/pkg"
	"main/expenses/pkg/validate"
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
	if good := validate.ID(id); !good {
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
	if good := validate.ID(ownerid); !good {
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
	var updateExpenseData globals.UpdateExpenseData
	err := globals.DecodeJSONBody(w, r, &updateExpenseData)
	if err != nil {
		var malformedreq *globals.MalformedReq
		if errors.As(err, &malformedreq) {
			http.Error(w, malformedreq.Msg, malformedreq.StatusCode)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

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
	if good := validate.All(id, name, spent); !good {
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
	var newExpenseData globals.NewExpenseData
	err := globals.DecodeJSONBody(w, r, &newExpenseData)
	if err != nil {
		var malformedreq *globals.MalformedReq
		if errors.As(err, &malformedreq) {
			http.Error(w, malformedreq.Msg, malformedreq.StatusCode)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	ownerid := r.URL.Query().Get("ownerid")
	name := r.URL.Query().Get("name")
	spent, err := strconv.Atoi(r.URL.Query().Get("spent"))
	if err != nil {
		var malformedreq *globals.MalformedReq
		if errors.As(err, &malformedreq) {
			http.Error(w, malformedreq.Msg, malformedreq.StatusCode)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	good := validate.All(ownerid, name, spent)
	if !good {
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
	if good := validate.ID(id); !good {
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

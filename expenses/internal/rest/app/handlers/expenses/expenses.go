package expenses

import (
	"errors"
	"main/expenses/internal/db/app"
	helpers "main/expenses/internal/rest/pkg"
	"main/expenses/pkg/validate"
	"net/http"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)
var log *zap.Logger
func SetLogger(logger *zap.Logger) {
	log = logger
}
func ErrInvalidFormat(w http.ResponseWriter) {
	http.Error(w, http.StatusText(400), 400)
	w.WriteHeader(400)
}
func ErrNotFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(404), 404)
	w.WriteHeader(404)
}
func GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if good := validate.ID(id); !good {
		ErrInvalidFormat(w)
		return
	}
	expense, err := db.GetExpenseById(id)
	if err != nil {
		ErrNotFound(w)
		return
	}
	err = helpers.WriteJSON(w, http.StatusOK, helpers.DefaultEnvelope(expense), nil)
	if err != nil {
		log.Error("Error writing JSON",
			zap.Error(err),
		)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
func GetByOwnerID(w http.ResponseWriter, r *http.Request) {
	ownerid := chi.URLParam(r, "id")
	if good := validate.ID(ownerid); !good {
		ErrInvalidFormat(w)
		return
	}
	ownerExpenses, err := db.GetExpensesByOwnerId(ownerid)
	if err != nil {
		ErrNotFound(w)
		return
	}
	err = helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"expenses": ownerExpenses}, nil)
	if err != nil {
		log.Error("Error writing JSON",
			zap.Error(err),
		)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	var updateExpenseData helpers.UpdateExpenseData
	err := helpers.DecodeJSONBody(w, r, &updateExpenseData)
	if err != nil {
		var malformedreq *helpers.MalformedReq
		if errors.As(err, &malformedreq) {
			http.Error(w, malformedreq.Msg, malformedreq.StatusCode)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	id := chi.URLParam(r, "id")
	name := updateExpenseData.Name
	inputSpent := updateExpenseData.Spent
	original, err := db.GetExpenseById(id)
	var spent int
	if len(name) <= 0 {
		name = original.Name
	}
	if inputSpent <= 0 {
		spent = original.Spent
	} else {
		spent = inputSpent
		if err != nil {
			ErrInvalidFormat(w)
			return
		}
	}
	if good := validate.All(id, name, spent); !good {
		ErrInvalidFormat(w)
		return
	}
	updatedExpense, err := db.UpdateExpense(id, name, spent)
	if err != nil {
		ErrNotFound(w)
		return
	}
	err = helpers.WriteJSON(w, http.StatusOK, helpers.DefaultEnvelope(updatedExpense), nil)
	if err != nil {
		log.Error("Error writing JSON",
			zap.Error(err),
		)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
func NewExpense(w http.ResponseWriter, r *http.Request) {
	var newExpenseData helpers.NewExpenseData
	err := helpers.DecodeJSONBody(w, r, &newExpenseData)
	if err != nil {
		var malformedreq *helpers.MalformedReq
		if errors.As(err, &malformedreq) {
			http.Error(w, malformedreq.Msg, malformedreq.StatusCode)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	ownerid := newExpenseData.OwnerID
	name := newExpenseData.Name
	spent := newExpenseData.Spent
	if err != nil {
		var malformedreq *helpers.MalformedReq
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
	newExpense, err := db.NewExpense(ownerid, name, spent)
	if err != nil {
		ErrNotFound(w)
		return
	}
	err = helpers.WriteJSON(w, http.StatusCreated, helpers.DefaultEnvelope(newExpense), nil)
	if err != nil {
		log.Error("Error writing JSON",
			zap.Error(err),
		)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
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

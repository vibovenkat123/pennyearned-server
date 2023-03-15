package expenses

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"main/expenses/internal/db/app"
	. "main/expenses/internal/rest/pkg"
	"main/expenses/pkg/validate"
	"net/http"
)

func GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if good := validate.ID(id); !good {
		App.WrongFormatResponse(w, r)
		return
	}
	expense, err := db.GetExpenseById(id)
	if err != nil {
		App.NotFoundResponse(w, r)
		return
	}
	err = App.WriteJSON(w, http.StatusOK, App.DefaultEnvelope(expense), nil)
	if err != nil {
		App.ServerErrorResponse(w, r, err)
	}
}
func GetByOwnerID(w http.ResponseWriter, r *http.Request) {
	ownerid := chi.URLParam(r, "id")
	if good := validate.ID(ownerid); !good {
		App.WrongFormatResponse(w, r)
		return
	}
	ownerExpenses, err := db.GetExpensesByOwnerId(ownerid)
	if err != nil {
		App.NotFoundResponse(w, r)
		return
	}
	err = App.WriteJSON(w, http.StatusOK, Envelope{"expenses": ownerExpenses}, nil)
	if err != nil {
		App.ServerErrorResponse(w, r, err)
		return
	}
}
func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	var updateExpenseData UpdateExpenseData
	err := App.DecodeJSONBody(w, r, &updateExpenseData)
	if err != nil {
		var malformedreq *MalformedReq
		if errors.As(err, &malformedreq) {
			App.ErrorResponse(w, r, malformedreq.StatusCode, malformedreq.Msg)
		} else {
			App.ServerErrorResponse(w, r, err)
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
			App.WrongFormatResponse(w, r)
			return
		}
	}
	if good := validate.All(id, name, spent); !good {
		App.WrongFormatResponse(w, r)
		return
	}
	updatedExpense, err := db.UpdateExpense(id, name, spent)
	if err != nil {
		App.NotFoundResponse(w, r)
		return
	}
	err = App.WriteJSON(w, http.StatusOK, App.DefaultEnvelope(updatedExpense), nil)
	if err != nil {
		App.ServerErrorResponse(w, r, err)
	}
}
func NewExpense(w http.ResponseWriter, r *http.Request) {
	var newExpenseData NewExpenseData
	err := App.DecodeJSONBody(w, r, &newExpenseData)
	if err != nil {
		var malformedreq *MalformedReq
		if errors.As(err, &malformedreq) {
			App.ErrorResponse(w, r, malformedreq.StatusCode, malformedreq.Msg)
		} else {
			App.ServerErrorResponse(w, r, err)
		}
		return
	}

	ownerid := newExpenseData.OwnerID
	name := newExpenseData.Name
	spent := newExpenseData.Spent
	if err != nil {
		var malformedreq *MalformedReq
		if errors.As(err, &malformedreq) {
			App.ErrorResponse(w, r, malformedreq.StatusCode, malformedreq.Msg)
		} else {
			App.ServerErrorResponse(w, r, err)
		}
		return
	}
	good := validate.All(ownerid, name, spent)
	if !good {
		App.WrongFormatResponse(w, r)
		return
	}
	newExpense, err := db.NewExpense(ownerid, name, spent)
	if err != nil {
		App.NotFoundResponse(w, r)
		return
	}
	err = App.WriteJSON(w, http.StatusCreated, App.DefaultEnvelope(newExpense), nil)
	if err != nil {
		App.ServerErrorResponse(w, r, err)
	}
}
func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if good := validate.ID(id); !good {
		App.WrongFormatResponse(w, r)
		return
	}
	err := db.DeleteExpense(id)
	if err != nil {
		App.NotFoundResponse(w, r)
		return
	}
	w.WriteHeader(204)
}

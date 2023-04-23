package expenses

import (
	"main/expenses/internal/db/app"
	dbHelpers "main/expenses/internal/db/pkg"
	. "main/expenses/internal/rest/pkg"
	"main/expenses/pkg/validate"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func GetByID(w http.ResponseWriter, r *http.Request) {
	App.Log.Info("Got request",
		zap.String("IP", r.RemoteAddr),
	)
	id := chi.URLParam(r, "id")
	if good := validate.ID(id); !good {
		App.BadRequestResponse(w, r, ErrInvalidID)
		return
	}
	expense, err := db.GetExpenseById(id)
	if err != nil {
		App.NotFoundResponse(w, r)
		return
	}
	err = App.WriteJSON(w, http.StatusOK, App.DefaultEnvelope(expense), nil)
	if err != nil {
		App.Log.Error("Got error while writing json",
			zap.Error(err),
		)
		App.ServerErrorResponse(w, r, err)
	}
}
func GetByOwnerID(w http.ResponseWriter, r *http.Request) {
	App.Log.Info("Got request",
		zap.String("IP", r.RemoteAddr),
	)
	ownerid := chi.URLParam(r, "id")
	if good := validate.ID(ownerid); !good {
		App.BadRequestResponse(w, r, ErrInvalidID)
		return
	}
	ownerExpenses, err := db.GetExpensesByOwnerId(ownerid)
	if err != nil {
		App.NotFoundResponse(w, r)
		return
	}
	err = App.WriteJSON(w, http.StatusOK, Envelope{"expenses": ownerExpenses}, nil)
	if err != nil {
		App.Log.Error("Got error while writing json",
			zap.Error(err),
		)
		App.ServerErrorResponse(w, r, err)
		return
	}
}
func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	App.Log.Info("Got request",
		zap.String("IP", r.RemoteAddr),
	)
	var updateExpenseData UpdateExpenseData
	err := App.ReadJSON(w, r, &updateExpenseData)
	if err != nil {
		App.BadRequestResponse(w, r, err)
		return
	}

	id := chi.URLParam(r, "id")
	name := updateExpenseData.Name
	spent := updateExpenseData.Spent
	expense, err := db.GetExpenseById(id)

	if err != nil {
		if err == dbHelpers.ErrExpenseNotFound {
			App.NotFoundResponse(w, r)
		} else {
			App.Log.Error("Got error while getting expense",
				zap.Error(err),
			)
			App.ServerErrorResponse(w, r, err)
		}
		return
	}

	if name != nil {
		expense.Name = *name
	}
	if spent != nil {
		expense.Spent = *spent
	}
	if good := validate.All(id, expense.Name, expense.Spent); !good {
		App.BadRequestResponse(w, r, ErrExpenseInvalid)
		return
	}
	updatedExpense, err := db.UpdateExpense(id, expense.Name, expense.Spent)
	if err != nil {
		App.NotFoundResponse(w, r)
		return
	}
	err = App.WriteJSON(w, http.StatusOK, App.DefaultEnvelope(updatedExpense), nil)
	if err != nil {
		App.Log.Error("Got error while writing json",
			zap.Error(err),
		)
		App.ServerErrorResponse(w, r, err)
	}
}
func NewExpense(w http.ResponseWriter, r *http.Request) {
	App.Log.Info("Got request",
		zap.String("IP", r.RemoteAddr),
	)
	var newExpenseData NewExpenseData
	err := App.ReadJSON(w, r, &newExpenseData)
	if err != nil {
		App.BadRequestResponse(w, r, err)
		return
	}
	ownerid := newExpenseData.OwnerID
	name := newExpenseData.Name
	spent := newExpenseData.Spent
	if err != nil {
		App.BadRequestResponse(w, r, err)
		return
	}
	good := validate.All(ownerid, name, spent)
	if !good {
		App.BadRequestResponse(w, r, ErrExpenseInvalid)
		return
	}
	newExpense, err := db.NewExpense(ownerid, name, spent)
	if err != nil {
		App.NotFoundResponse(w, r)
		return
	}
	err = App.WriteJSON(w, http.StatusCreated, App.DefaultEnvelope(newExpense), nil)
	if err != nil {
		App.Log.Error("Got error while writing json",
			zap.Error(err),
		)
		App.ServerErrorResponse(w, r, err)
	}
}
func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	App.Log.Info("Got request",
		zap.String("IP", r.RemoteAddr),
	)
	id := chi.URLParam(r, "id")
	if good := validate.ID(id); !good {
		App.BadRequestResponse(w, r, ErrInvalidID)
		return
	}
	err := db.DeleteExpense(id)
	if err != nil {
		App.NotFoundResponse(w, r)
		return
	}
	w.WriteHeader(204)
}

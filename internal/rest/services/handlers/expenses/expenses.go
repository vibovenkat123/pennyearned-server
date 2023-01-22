package expensesHandler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/vibovenkat123/pennyearned-server/internal/db/helpers"
	"net/http"
)

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
func ExpensesCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var expense []dbHelpers.Expenses
		if uID := chi.URLParam(r, "ownerId"); uID != "" {
			expense = dbHelpers.GetExpensesByOwnerId(uID)
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), "expenses", expense)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func GetExpenses(w http.ResponseWriter, r *http.Request) {
	expenses := r.Context().Value("expenses").([]dbHelpers.Expenses)
	render.JSON(w, r, expenses)
}

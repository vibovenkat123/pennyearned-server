package expensesHandler


import (
	"github.com/go-chi/chi/v5"
	"net/http"
    "context"
    "github.com/vibovenkat123/pennyearned-server/internal/db/helpers"
    "github.com/go-chi/render"
)
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
func ExpensesCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var expense []dbHelpers.Expenses
		if eID := chi.URLParam(r, "expensesID"); eID != "" {
            expense = dbHelpers.GetExpensesByOwnerId(eID) 
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

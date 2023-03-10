package api

import (
	"fmt"
	"main/expenses/internal/rest/app/handlers/expenses"
	helpers "main/expenses/internal/rest/pkg"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

var env string

func Expose(log *zap.Logger) {
    if helpers.ConvertErr != nil {
        log.Error("The port variable is not a valid int",
            zap.Error(helpers.ConvertErr),
        )
    }
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/api/expense", expensesRouter)
	r.Route("/api/users", userRouter)
	fmt.Printf("Starting %v server on port :%v\n", env, helpers.Port)
	handler := cors.Default().Handler(r)
	http.ListenAndServe(fmt.Sprintf(":%v", helpers.Port), handler)
}
func userRouter(r chi.Router) {
	r.Route("/{id}", UserIDRouter)
}
func UserIDRouter(r chi.Router) {
	r.Get("/expenses", expenses.GetByOwnerID)
}
func expensesRouter(r chi.Router) {
	r.Route("/{id}", ExpenseIDRouter)
	r.Post("/", expenses.NewExpense)
}
func ExpenseIDRouter(r chi.Router) {
	r.Get("/", expenses.GetByID)
	r.Delete("/", expenses.DeleteExpense)
	r.Patch("/", expenses.UpdateExpense)
}

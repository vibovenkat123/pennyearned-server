package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"os"
	"time"
    expenseHandler "github.com/vibovenkat123/pennyearned-server/internal/rest/services/handlers/expenses"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
// function to expose all the api routes
func Expose() {
	port := 4000
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/api/expenses", func(r chi.Router) {
		r.Route("/{expensesID}", func(r chi.Router) {
			r.Use(expenseHandler.ExpensesCtx)
            r.Get("/", expenseHandler.GetExpenses)
		})
	})
    fmt.Printf("listening on port %v\n", port)
    http.ListenAndServe(fmt.Sprintf(":%v", port), r)
}

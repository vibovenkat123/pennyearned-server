package api

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"go.uber.org/zap"
	handlers "main/expenses/internal/rest/app/handlers/expenses"
	. "main/expenses/internal/rest/pkg"
	"net/http"
	"time"
)

var env string
var adapter *chiadapter.ChiLambda

func Expose(local bool) {
	r := chi.NewRouter()
	r.NotFound(http.HandlerFunc(App.NotFoundResponse))
	r.MethodNotAllowed(http.HandlerFunc(App.MethodNotAllowedResponse))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/v1/api/expense", expenseRouter)
	r.Route("/v1/api/user", userRouter)
	if local {
		handler := cors.Default().Handler(r)
		App.Log.Info("Starting server",
			zap.String("Port", fmt.Sprintf("%v", App.Conf.Port)),
		)
		http.ListenAndServe(fmt.Sprintf(":%v", App.Conf.Port), handler)
	} else {
		App.Log.Info("Starting server on lambda")
		adapter = chiadapter.New(r)
	}
}
func StartAPI() {
	Expose(Local)
	if !Local {
		lambda.Start(Handler)
	}
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return adapter.ProxyWithContext(ctx, req)
}

func userRouter(r chi.Router) {
	r.Route("/{id}", UserIDRouter)
}
func UserIDRouter(r chi.Router) {
	r.Get("/expenses", handlers.GetByOwnerID)
}
func expenseRouter(r chi.Router) {
	r.Route("/{id}", ExpenseIDRouter)
	r.Post("/", handlers.NewExpense)
}
func ExpenseIDRouter(r chi.Router) {
	r.Get("/", handlers.GetByID)
	r.Delete("/", handlers.DeleteExpense)
	r.Patch("/", handlers.UpdateExpense)
}

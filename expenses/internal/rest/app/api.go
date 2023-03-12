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
	"main/expenses/internal/rest/app/handlers/expenses"
	helpers "main/expenses/internal/rest/pkg"
	"net/http"
	"time"
)

var env string
var adapter *chiadapter.ChiLambda

func Expose(log *zap.Logger, local bool) {
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
	r.Route("/api/user", userRouter)
	if local {
		handler := cors.Default().Handler(r)
		log.Info("Starting server",
			zap.String("Port", fmt.Sprintf("%v", helpers.Port)),
		)
		http.ListenAndServe(fmt.Sprintf(":%v", helpers.Port), handler)
	} else {
		log.Info("Starting server on lambda")
		adapter = chiadapter.New(r)
	}
}
func StartAPI(log *zap.Logger) {
	Expose(log, helpers.Local)
	if !helpers.Local {
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

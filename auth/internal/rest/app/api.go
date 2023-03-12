package api

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	users "main/auth/internal/rest/app/handlers/users"
	helpers "main/auth/internal/rest/pkg"
	"time"
)

var env string
var adapter *chiadapter.ChiLambda

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
	r.Route("/api/user", userRouter)
	adapter = chiadapter.New(r)
}
func StartAPI(log *zap.Logger) {
	Expose(log)
	log.Info("Starting server",
		zap.String("Port", fmt.Sprintf("%v", helpers.Port)),
	)
	lambda.Start(Handler)
}
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return adapter.ProxyWithContext(ctx, req)
}
func userRouter(r chi.Router) {
	r.Post("/session", users.SignIn)
	r.Post("/", users.SignUp)
	r.Delete("/session/{id}", users.SignOut)
	r.Get("/{id}", users.GetByCookie)
}

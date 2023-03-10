package api

import (
	"fmt"
	helpers "main/auth/internal/rest/pkg"
	users "main/auth/internal/rest/app/handlers/users"
	"net/http"
	"time"
    "log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)
var env string
func Expose() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/api/user", userRouter)
	fmt.Printf("Starting %v server on port :%v\n", env, helpers.Port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%v", helpers.Port), r))
}
func userRouter(r chi.Router) {
	r.Post("/session", users.SignIn)
	r.Post("/", users.SignUp)
	r.Delete("/session/{id}", users.SignOut)
	r.Get("/{id}", users.GetByCookie)
}

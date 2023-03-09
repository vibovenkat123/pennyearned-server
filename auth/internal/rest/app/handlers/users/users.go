package users

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"main/auth/internal/db/app"
	helpers "main/auth/internal/rest/pkg"
	"main/auth/pkg/userFunctions"
	"net/http"
)

func ErrInvalidFormat(w http.ResponseWriter, text string) {
	http.Error(w, text, 400)
	w.WriteHeader(400)
}
func ErrNotFound(w http.ResponseWriter, text string) {
	http.Error(w, text, 404)
	w.WriteHeader(404)
}
func ErrAlreadyFound(w http.ResponseWriter, text string) {
	http.Error(w, text, 409)
	w.WriteHeader(409)
}
func Success(w http.ResponseWriter, text string) {
	http.Error(w, text, 200)
	w.WriteHeader(200)
}
func SuccessfullyDeleted(w http.ResponseWriter, text string) {
	http.Error(w, text, 204)
	w.WriteHeader(204)
}
func SignIn(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")
	if !(userFunctions.PasswordCheck(password) && userFunctions.EmailCheck(email)) {
		ErrNotFound(w, http.StatusText(404))
		return
	}
	cookie, err := db.SignIn(email, password, r.Context())
	if err != nil {
		ErrNotFound(w, http.StatusText(404))
		return
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(200)
}
func SignUp(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")
	name := r.URL.Query().Get("name")
	username := r.URL.Query().Get("username")
	if !userFunctions.Validate(name, username, email, password) {
		ErrInvalidFormat(w, http.StatusText(400))
		return
	}
	cookie, err := db.SignUp(name, username, email, password, r.Context())
	if err != nil {
		ErrAlreadyFound(w, http.StatusText(409))
		return
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(201)
}
func GetByCookie(w http.ResponseWriter, r *http.Request) {
	cookieID := chi.URLParam(r, "id")
	encoder := json.NewEncoder(w)
	val, err := db.GetByCookie(cookieID, r.Context())
	if err != nil {
		ErrNotFound(w, http.StatusText(404))
		return
	}
	res := helpers.UserIDRes{
		ID: val,
	}
	encoder.Encode(res)
	w.WriteHeader(200)
}
func SignOut(w http.ResponseWriter, r *http.Request) {
	cookieID := chi.URLParam(r, "id")
	cookie, err := db.SignOut(cookieID, r.Context())
	if err != nil {
		ErrNotFound(w, http.StatusText(404))
		return
	}
	http.SetCookie(w, cookie)
	SuccessfullyDeleted(w, http.StatusText(204))
}

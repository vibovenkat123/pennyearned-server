package users

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"main/auth/internal/db/app"
	dbGlobals "main/auth/internal/db/pkg"
	globals "main/auth/internal/rest/pkg"
	"main/auth/pkg/validate"
	"net/http"
)

func ErrServer(w http.ResponseWriter, text string) {
	http.Error(w, text, 500)
	w.WriteHeader(500)
}
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
	w.WriteHeader(200)
	w.Write([]byte(text))
}
func PartialSuccess(w http.ResponseWriter, text string) {
	w.WriteHeader(202)
	w.Write([]byte(text))
}
func SuccessfullyDeleted(w http.ResponseWriter, text string) {
	w.WriteHeader(204)
	w.Write([]byte(text))
}
func SignIn(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")
	if !validate.Password(password) {
		ErrNotFound(w, http.StatusText(404))
		return
	}
	accessToken, err := db.SignIn(email, password, r.Context())
	if err != nil {
		ErrNotFound(w, http.StatusText(404))
		return
	}
	response := globals.UserAccessRes{
		AccessToken: *accessToken,
	}
	encoder.Encode(response)
	w.WriteHeader(200)
}
func SendVerification(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	to := []string{email}
	err := db.SendEmail(to, r.Context())
	if err != nil {
		ErrServer(w, err.Error())
		return
	}
	PartialSuccess(w, http.StatusText(202))
}
func SignUp(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	code := chi.URLParam(r, "code")
	password := r.URL.Query().Get("password")
	name := r.URL.Query().Get("name")
	username := r.URL.Query().Get("username")
	if !validate.All(name, username, code, password) {
		ErrInvalidFormat(w, http.StatusText(400))
		return
	}
	accessToken, err := db.SignUp(name, username, password, code, r.Context())
	if err != nil {
		if err == dbGlobals.ErrInvalidCode {
			ErrNotFound(w, http.StatusText(404))
			return
		}
		ErrAlreadyFound(w, http.StatusText(409))
		return
	}
	response := globals.UserAccessRes{
		AccessToken: *accessToken,
	}
	encoder.Encode(response)
	w.WriteHeader(201)
}
func GetByCookie(w http.ResponseWriter, r *http.Request) {
	cookieID := chi.URLParam(r, "id")
	encoder := json.NewEncoder(w)
	val, err := db.GetByAccess(cookieID, r.Context())
	if err != nil {
		ErrNotFound(w, http.StatusText(404))
		return
	}
	res := globals.UserIDRes{
		ID: val,
	}
	encoder.Encode(res)
	w.WriteHeader(200)
}
func SignOut(w http.ResponseWriter, r *http.Request) {
	accessToken := chi.URLParam(r, "id")
	err := db.SignOut(accessToken, r.Context())
	if err != nil {
		ErrNotFound(w, http.StatusText(404))
		return
	}
	SuccessfullyDeleted(w, http.StatusText(204))
}

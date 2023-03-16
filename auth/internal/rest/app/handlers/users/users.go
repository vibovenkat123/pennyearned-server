package users

import (
	"main/auth/internal/db/app"
	dbGlobals "main/auth/internal/db/pkg"
	. "main/auth/internal/rest/pkg"
	"main/auth/pkg/validate"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

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
	var signInData SignInData
	err := App.ReadJSON(w, r, &signInData)
	if err != nil {
		App.BadRequestResponse(w, r, err)
		return
	}
	email := signInData.Email
	password := signInData.Password
	if !validate.Password(password) {
		App.NotFoundResponse(w, r)
		return
	}
	accessToken, err := db.SignIn(email, password, r.Context())
	if err != nil {
		App.NotFoundResponse(w, r)
		return
	}
	accessTokenResponse := UserAccessRes{
		AccessToken: *accessToken,
	}
	err = App.WriteJSON(w, http.StatusOK, App.DefaultEnvelope(accessTokenResponse), nil)
	if err != nil {
		App.ServerErrorResponse(w, r, err)
	}
}

func SendVerification(w http.ResponseWriter, r *http.Request) {
	var verifyData VerifyData
	err := App.ReadJSON(w, r, &verifyData)
	if err != nil {
		App.ServerErrorResponse(w, r, err)
		return
	}
	email := verifyData.Email
	to := []string{email}
	err = db.SendEmail(to, r.Context())
	if err != nil {
		if strings.Contains(err.Error(), "is not a valid RFC-5321 address") {
			App.BadRequestResponse(w, r, ErrEmailWrongFormat)
			return
		}
		App.ServerErrorResponse(w, r, err)
		return
	}
	PartialSuccess(w, http.StatusText(202))
}
func SignUpVerify(w http.ResponseWriter, r *http.Request) {
	var signUpVerifyData SignUpVerifyData
	err := App.ReadJSON(w, r, &signUpVerifyData)
	if err != nil {
		App.BadRequestResponse(w, r, err)
		return
	}
	code := chi.URLParam(r, "code")
	password := signUpVerifyData.Password
	name := signUpVerifyData.Name
	username := signUpVerifyData.Username
	if !validate.All(name, username, code, password) {
		App.BadRequestResponse(w, r, ErrSignUPWrongFormat)
		return
	}
	accessToken, err := db.SignUp(name, username, password, code, r.Context())
	if err != nil {
		if err == dbGlobals.ErrInvalidCode {
			App.NotFoundResponse(w, r)
			return
		}
		App.ConflictResponse(w, r)
		return
	}
	accessTokenResponse := UserAccessRes{
		AccessToken: *accessToken,
	}
	err = App.WriteJSON(w, http.StatusCreated, App.DefaultEnvelope(accessTokenResponse), nil)
	if err != nil {
		App.ServerErrorResponse(w, r, err)
		return
	}
}
func GetByCookie(w http.ResponseWriter, r *http.Request) {
	cookieID := chi.URLParam(r, "id")
	val, err := db.GetByAccess(cookieID, r.Context())
	if err != nil {
		App.NotFoundResponse(w, r)
		return
	}
	userIDResponse := UserIDRes{
		ID: val,
	}
	err = App.WriteJSON(w, http.StatusOK, App.DefaultEnvelope(userIDResponse), nil)
	if err != nil {
		App.ServerErrorResponse(w, r, err)
		return
	}
}
func SignOut(w http.ResponseWriter, r *http.Request) {
	accessToken := chi.URLParam(r, "id")
	err := db.SignOut(accessToken, r.Context())
	if err != nil {
		App.NotFoundResponse(w, r)
		return
	}
	SuccessfullyDeleted(w, http.StatusText(204))
}

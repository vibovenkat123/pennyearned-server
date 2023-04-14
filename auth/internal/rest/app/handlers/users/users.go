package users

import (
	"main/auth/internal/db/app"
	dbGlobals "main/auth/internal/db/pkg"
	. "main/auth/internal/rest/pkg"
	"main/auth/pkg/validate"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func PartialSuccess(w http.ResponseWriter, text string, r *http.Request) {
	w.WriteHeader(202)
	_, err := w.Write([]byte(text))
	if err != nil {
		App.LogError(err, r)
	}
}
func SuccessfullyDeleted(w http.ResponseWriter, text string, r *http.Request) {
	w.WriteHeader(204)
	_, err := w.Write([]byte(text))
	if err != nil {
		App.LogError(err, r)
	}
}
func SignIn(w http.ResponseWriter, r *http.Request) {
	App.Log.Info("Got request",
		zap.String("IP", r.RemoteAddr),
	)
	App.Log.Info("Signing in");
	var signInData SignInData
	App.Log.Info("Reading JSON")
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
		App.Log.Info("The user was not found")
		App.NotFoundResponse(w, r)
		return
	}
	accessTokenResponse := UserAccessRes{
		AccessToken: *accessToken,
	}
	err = App.WriteJSON(w, http.StatusOK,
		App.DefaultEnvelope(accessTokenResponse), nil)
	if err != nil {
		App.Log.Error("Error while signing in",
			zap.Error(err),
		)
		App.ServerErrorResponse(w, r, err)
	}
}

func SendVerification(w http.ResponseWriter, r *http.Request) {
	App.Log.Info("Got request",
		zap.String("IP", r.RemoteAddr),
	)
	App.Log.Info("Sending verification")
	var verifyData VerifyData
	App.Log.Info("Reading JSON")
	err := App.ReadJSON(w, r, &verifyData)
	if err != nil {
		App.ServerErrorResponse(w, r, err)
		return
	}
	email := verifyData.Email
	to := []string{email}
	App.Log.Info("Sending email")
	err = db.SendEmail(to, r.Context())
	if err != nil {
		if strings.Contains(err.Error(), "is not a valid RFC-5321 address") {
			App.Log.Info("The email is not valid")
			App.BadRequestResponse(w, r, ErrEmailWrongFormat)
			return
		}
		App.Log.Error("Error while sending verification",
			zap.Error(err),
		)
		App.ServerErrorResponse(w, r, err)
		return
	}
	PartialSuccess(w, http.StatusText(202), r)
}
func SignUpVerify(w http.ResponseWriter, r *http.Request) {
	App.Log.Info("Got request",
		zap.String("IP", r.RemoteAddr),
	)
	App.Log.Info("Signing up and verifying")
	var signUpVerifyData SignUpVerifyData
	App.Log.Info("Reading JSON")
	err := App.ReadJSON(w, r, &signUpVerifyData)
	if err != nil {
		App.BadRequestResponse(w, r, err)
		return
	}
	App.Log.Info("Reading verify code")
	code := chi.URLParam(r, "code")
	password := signUpVerifyData.Password
	username := signUpVerifyData.Username
	App.Log.Info("Validating code, username, and pass")
	if !validate.All(username, code, password) {
		App.BadRequestResponse(w, r, ErrSignUPWrongFormat)
		return
	}
	App.Log.Info("Signing up")
	accessToken, err := db.SignUp(username, password, code, r.Context())
	App.Log.Debug("Signed Up",
		zap.String("Access token", *accessToken),
	)
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
		App.Log.Error("Error while signing up",
			zap.Error(err),
		)
		App.ServerErrorResponse(w, r, err)
		return
	}
}
func GetByCookie(w http.ResponseWriter, r *http.Request) {
	App.Log.Info("Got request",
		zap.String("IP", r.RemoteAddr),
	)
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
		App.Log.Error("Cannot get cookie",
			zap.Error(err),
		)
		App.ServerErrorResponse(w, r, err)
		return
	}
}
func SignOut(w http.ResponseWriter, r *http.Request) {
	App.Log.Info("Got request",
		zap.String("IP", r.RemoteAddr),
	)
	accessToken := chi.URLParam(r, "id")
	err := db.SignOut(accessToken, r.Context())
	if err != nil {
		App.NotFoundResponse(w, r)
		return
	}
	SuccessfullyDeleted(w, http.StatusText(204), r)
}

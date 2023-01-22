package expensesHandler

import (
	"github.com/vibovenkat123/pennyearned-server/internal/db/helpers"
	"net/http"
)

// structs
type Handler func(w http.ResponseWriter, r *http.Request) error
type Response struct {
	ok bool
}
type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}
type UserPayload struct {
	*dbHelpers.User
}
type ExpensesResponse struct {
	*dbHelpers.Expenses
	User *UserPayload
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

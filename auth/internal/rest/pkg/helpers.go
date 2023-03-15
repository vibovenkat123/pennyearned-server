package apiHelpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type MalformedReq struct {
	StatusCode int
	Msg        string
}

// a basic error log
func (malformedreq *MalformedReq) Error() string {
	return malformedreq.Msg
}

func (app *Application) ReadJSON(w http.ResponseWriter, r *http.Request, dataToDecode interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	err := json.NewDecoder(r.Body).Decode(dataToDecode)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	return nil
}

func (app *Application) DefaultEnvelope(data any) Envelope {
	return Envelope{topKey: data}
}

func (app *Application) WriteJSON(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	// Encode the data to JSON, returning the error if there was one.
	// Make the json pretty printed/tabbed
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}
	// Add the "Content-Type: application/json" header, then write the status code and // JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

package api
import (
    "net/http"
    "main/expenses/internal/rest/services/handlers/expenses"
    helpers "main/expenses/internal/rest/helpers"
    "fmt"
    "log"
    "os"
)
func Expose() {
    logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)
    mux := http.NewServeMux()
    mux.HandleFunc(helpers.GetExpenseByIDUrl, expenses.GetByID)
    mux.HandleFunc(helpers.GetExpensesByOwnerIDUrl, expenses.GetByOwnerID)
    srv :=  &http.Server{
        Addr: fmt.Sprintf(":%d", helpers.Port),
        Handler: mux,
    }
    logger.Printf("starting %s server on %s", os.Getenv("GO_ENV"), srv.Addr)
    err := srv.ListenAndServe()
    logger.Fatal(err)
}


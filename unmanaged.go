// +build !appengine

package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/juju/errors"
	"golang.org/x/net/context"
)

func init() {
	AddMiddleware(LoggerMiddleware)
	AddMiddleware(ClientErrorMiddleware)
}

// LoggerMiddleware logs to stderr the error.
func LoggerMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error {
	if err := next(ctx); err != nil {
		log.Printf("HANDLER ERROR\n%s\n", errors.ErrorStack(err))
	}

	return nil
}

// ClientErrorMiddleware rejects the request with a 500 status code when an error occurs.
func ClientErrorMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error {
	err := next(ctx)
	if err != nil {
		w.(*appResponseWriter).Reset()
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if os.Getenv("DEBUG") == "true" {
			http.Error(w, errors.ErrorStack(err), http.StatusInternalServerError)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	return err
}

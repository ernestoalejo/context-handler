// +build !appengine

package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/juju/errors"
	"golang.org/x/net/context"
)

func init() {
	AddMiddleware(LoggerMiddleware)
	AddMiddleware(ClientErrorMiddleware)
}

// LoggerMiddlewares logs to stderr the error.
func LoggerMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error {
	if err := next(); err != nil {
		separator := "--------------------------------------------------"
		fmt.Fprintf(os.Stderr, "HANDLER ERROR:\n%s\n%s\n%s\n", separator, errors.ErrorStack(err), separator)
	}

	return nil
}

// ClientErrorMiddleware rejects the request with a 500 status code when an error occurs.
func ClientErrorMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error {
	err := next()
	if err != nil {
		http.Error(w, "handler error", http.StatusInternalServerError)
	}

	return err
}

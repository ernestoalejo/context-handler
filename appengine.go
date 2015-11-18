// +build appengine

package handler

import (
	"net/http"

	"github.com/juju/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

func init() {
	AddMiddleware(LoggerMiddleware)
	AddMiddleware(ClientErrorMiddleware)
}

// LoggerMiddleware logs to Cloud Logging the error.
func LoggerMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error {
	if err := next(); err != nil {
		log.Errorf(ctx, "%s", errors.ErrorStack(err))
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

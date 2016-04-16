// +build appengine

package handler

import (
	"net/http"
	"time"

	"github.com/juju/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	AddMiddleware(ContextMiddleware)
	AddMiddleware(ContextTimeout)
	AddMiddleware(LoggerMiddleware)
	AddMiddleware(ClientErrorMiddleware)
}

// ContextMiddleware bootstraps the context to be able to contact with App Engine.
func ContextMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error {
	return next(appengine.WithContext(ctx, r))
}

// ContextTimeout sets a timeout in the context.
func ContextTimeout(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error {
	secs := 55 * time.Second
	if r.Header.Get("X-AppEngine-QueueName") != "" || r.Header.Get("X-AppEngine-Cron") != "" {
		secs = (9*60 + 55) * time.Second
	}

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(secs))

	done := make(chan error, 1)
	go func() {
		done <- next(ctx)
	}()

	select {
	case err := <-done:
		cancel()
		return err

	case <-time.After(secs):
		log.Criticalf(ctx, "request timeout")
		cancel()
		return nil
	}

	panic("should not reach here")
}

// ContextTimeoutFailFast sets a shorter timeout in the context of interactive requests and allows
// the previous middleware to log or answer request timeouts. Be sure you process the
// "request timeout" error or it will panic when finishing the middleware stack. It is not activated
// by default; but you can clear the middleware stack and put it instead of ContextTimeout().
func ContextTimeoutFailFast(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error {
	secs := 15 * time.Second
	if r.Header.Get("X-AppEngine-QueueName") != "" || r.Header.Get("X-AppEngine-Cron") != "" {
		secs = (9*60 + 55) * time.Second
	}

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(secs))

	done := make(chan error, 1)
	go func() {
		done <- next(ctx)
	}()

	select {
	case err := <-done:
		cancel()
		return err

	case <-time.After(secs):
		cancel()
		return errors.New("request timeout")
	}

	panic("should not reach here")
}

// LoggerMiddleware logs to Cloud Logging the error.
func LoggerMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error {
	if err := next(ctx); err != nil {
		log.Errorf(ctx, "%s", errors.ErrorStack(err))
	}

	return nil
}

// ClientErrorMiddleware rejects the request with a 500 status code when an error occurs.
func ClientErrorMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error {
	err := next(ctx)
	if err != nil {
		w.(*AppResponseWriter).Reset()
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if appengine.IsDevAppServer() {
			http.Error(w, errors.ErrorStack(err), http.StatusInternalServerError)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	return err
}

package handler

import (
	"net/http"

	"golang.org/x/net/context"
)

func init() {
	AddMiddleware(loggerMiddleware)
	AddMiddleware(clientErrorMiddleware)
}

// Middleware should be implemented by the functions that want to intercept requests
// or responses.
type Middleware func(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error

// NextMiddlewareFn should be called after the processing in the middleware finishes to follow the chain. An
// instance of this function type will be passed to each middleware call.
type NextMiddlewareFn func() error

var middlewares []Middleware

// AddMiddleware to the list of middlewares.
func AddMiddleware(middleware Middleware) {
	middlewares = append(middlewares, middleware)
}

func runMiddlewares(ctx context.Context, w http.ResponseWriter, r *http.Request, handler CtxHandler, current int) error {
	if current >= len(middlewares) {
		return handler(ctx, w, r)
	}

	err := middlewares[current](ctx, w, r, func() error {
		return runMiddlewares(ctx, w, r, handler, current+1)
	})
	return err
}

func clientErrorMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error {
	err := next()
	if err != nil {
		http.Error(w, "handler error", http.StatusInternalServerError)
	}

	return err
}

package handler

import (
	"net/http"

	"golang.org/x/net/context"
)

// Middleware should be implemented by the functions that want to intercept requests
// or responses.
type Middleware func(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error

// NextMiddlewareFn should be called after the processing in the middleware finishes to follow the chain. An
// instance of this function type will be passed to each middleware call.
type NextMiddlewareFn func(ctx context.Context) error

var middlewares []Middleware

// AddMiddleware to the list of middlewares.
func AddMiddleware(middleware Middleware) {
	middlewares = append(middlewares, middleware)
}

// ClearMiddlewares cleans the list of applied middlewares to let the app choose the order
// and what ones get activated with the requests.
func ClearMiddlewares() {
	middlewares = []Middleware{}
}

func runMiddlewares(ctx context.Context, w http.ResponseWriter, r *http.Request, handler CtxHandler, current int) error {
	if current >= len(middlewares) {
		return handler(ctx, w, r)
	}

	err := middlewares[current](ctx, w, r, func(ctx context.Context) error {
		return runMiddlewares(ctx, w, r, handler, current+1)
	})
	return err
}

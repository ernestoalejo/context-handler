// +build appengine

package handler

import (
	"net/http"

	"github.com/juju/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

func loggerMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error {
	if err := next(); err != nil {
		log.Errorf(ctx, "%s", errors.ErrorStack(err))
	}

	return nil
}

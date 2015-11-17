// +build !appengine

package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/juju/errors"
	"golang.org/x/net/context"
)

func loggerMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddlewareFn) error {
	if err := next(); err != nil {
		separator := "--------------------------------------------------"
		fmt.Fprintf(os.Stderr, "HANDLER ERROR:\n%s\n%s\n%s\n", separator, errors.ErrorStack(err), separator)
	}

	return nil
}

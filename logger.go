package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/juju/errors"
	"golang.org/x/net/context"
)

type Logger func(ctx context.Context, r *http.Request, err error)

var activeLogger Logger

func init() {
	SetActiveLogger(func(ctx context.Context, r *http.Request, err error) {
		separator := "--------------------------------------------------"
		fmt.Fprintf(os.Stderr, "HANDLER ERROR:\n%s\n%s\n%s\n", separator, errors.ErrorStack(err), separator)
	})
}

// SetActiveLogger changes the active logger for errors. By default it prints the error in stderr.
func SetActiveLogger(logger Logger) {
	activeLogger = logger
}

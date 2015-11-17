package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/juju/errors"
	"golang.org/x/net/context"
)

// CtxHandler should be implemented by the handlers.
type CtxHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// ReadJSON checks the request to see if it's a POST one; and reads the JSON data.
func ReadJSON(r *http.Request, data interface{}) error {
	if r.Method != "POST" {
		return errors.New("bad method")
	}

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return errors.Trace(err)
	}

	return nil
}

type handlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type appResponseWriter struct {
	statusCode int
	buffer     *bytes.Buffer
	header     http.Header
	written    bool
}

func (w *appResponseWriter) Header() http.Header {
	return w.header
}

func (w *appResponseWriter) Write(data []byte) (n int, err error) {
	if !w.written {
		w.statusCode = http.StatusOK
		w.written = true
	}

	return w.buffer.Write(data)
}

func (w *appResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.written = true
}

// Ctx adapts a context handler to the standard HTTP lib contract.
func Ctx(fn CtxHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wbuf := &appResponseWriter{
			buffer: bytes.NewBuffer(nil),
			header: make(http.Header),
		}

		ctx := context.Background()

		if err := fn(ctx, wbuf, r); err != nil {
			log.Println(errors.ErrorStack(err))
			http.Error(wbuf, "handler error", http.StatusInternalServerError)
		}

		for k, v := range wbuf.header {
			w.Header()[k] = v
		}
		w.WriteHeader(wbuf.statusCode)
		w.Write(wbuf.buffer.Bytes())
	}
}
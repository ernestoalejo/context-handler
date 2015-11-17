package handler

import (
	"bytes"
	"net/http"

	"golang.org/x/net/context"
)

// CtxHandler should be implemented by the handlers.
type CtxHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

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

		if err := runMiddlewares(ctx, wbuf, r, fn, 0); err != nil {
			http.Error(wbuf, "unhandled error", http.StatusInternalServerError)
			return
		}

		for k, v := range wbuf.header {
			w.Header()[k] = v
		}
		w.WriteHeader(wbuf.statusCode)
		w.Write(wbuf.buffer.Bytes())
	}
}

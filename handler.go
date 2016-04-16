package handler

import (
	"bytes"
	"net/http"

	"golang.org/x/net/context"
)

// CtxHandler should be implemented by the handlers.
type CtxHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type handlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// AppResponseWriter the package uses to cache the responses. Exposed to be able to build or own middlewares.
type AppResponseWriter struct {
	statusCode int
	buffer     *bytes.Buffer
	header     http.Header
	written    bool
}

// Header will call the same function in the standard response writer.
func (w *AppResponseWriter) Header() http.Header {
	return w.header
}

// Write will call the same function in the standard response writer.
func (w *AppResponseWriter) Write(data []byte) (n int, err error) {
	if !w.written {
		w.statusCode = http.StatusOK
		w.written = true
	}

	return w.buffer.Write(data)
}

// WriteHeader will call the same function in the standard response writer.
func (w *AppResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.written = true
}

// Reset removes any content in the memory buffer for the response.
func (w *AppResponseWriter) Reset() {
	w.buffer.Reset()
	w.header = make(http.Header)
}

// Ctx adapts a context handler to the standard HTTP lib contract.
func Ctx(fn CtxHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wbuf := &AppResponseWriter{
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

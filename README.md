
# context-handler

[![GoDoc](https://godoc.org/github.com/ernestoalejo/context-handler?status.svg)](https://godoc.org/github.com/ernestoalejo/context-handler)

Go HTTP handlers using context.


## Installation

```shell
go get github.com/ernestoalejo/context-handler
```


## Features

 - Pass context to all handlers.
 - Buffer responses to return errors to users correctly.
 - Read JSON in POST requests easily.
 - Return errors easily from the handler.
 - Show stack traces of errors.


### Usage

```go
package main

import (
  "fmt"
  "net/http"

  "github.com/ernestoalejo/context-handler"
  "github.com/juju/errors"
  "golang.org/x/net/context"
)

func main() {
  http.HandleFunc("/", handler.Ctx(homeHandler))
  http.HandleFunc("/health", handler.Ctx(healthHandler))
  http.ListenAndServe(":8080", nil)
}

func homeHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
  fmt.Fprintln(w, "Home handler")

  return errors.New("cannot handle request correctly")
}

func healthHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
  fmt.Fprintln(w, "ok")
  
  return nil
}
```

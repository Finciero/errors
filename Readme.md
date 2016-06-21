# Errors

Small library to manage errors across all services in the **Finciero API**.
Provides standard status codes declaration and convenience functions to
easily creates an error of the given status code, and wrapping an existing
error.

Also provides methods to pass this errors through **gRPC** server as a `grpcErr`.

> NOTICE: For now no documentation is available, but tests can be used as examples.

## Back end example usage

```go
package  server

import (
    "github.com/Finciero/errors"
)

type server struct {}

func (s *server) Req() (*Res, error) {
  res, err := fn()

  if err != nil {
    err = errors.New(StatusInternalServer, "unexpected error")
    return nil,  err.ToGRPC()
  }

  return res, nil
}
```

## Front end example usage

```go
package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "log"

    "github.com/Finciero/errors"
)

var (
  grpcClient Client
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  w.Header.Set("Content-Type", "application/json; charset=UTF-8")

  req, err := grpcClient.Req()
  if err != nil {
    err = errors.FromGRPC(err)
    switch err.StatusCode {
    case errors.StatusBadRequest:
      // manage bad requests
      w.WriteHeader(err.StatusCode)
      json.NewEncoder(w).Encoder(err) // there is no much more to do in case of failure
      return
    case err.StatusInternalServer:
      // manage internal server
      w.WriteHeader(err.StatusCode)
      json.NewEncoder(w).Encoder(err) // there is no much more to do in case of failure
      return
    default:
      // other unexpected error should be considered an internal server error
      err = errors.NewFromError(errors.StatusInternalServer, err)
      w.WriteHeader(err.StatusCode)
      json.NewEncoder(w).Encode(err) // at this point we need to ignore a possible error
      return
    }
  }

  if err := json.Decoder(w).Decode(err); err != nil {
      err = errors.NewFromError(errors.StatusInternalServer, err)
      w.WriteHeader(err.StatusCode)
      json.NewEncoder(w).Encode(err) // at this point we need to ignore a possible error
  }
}

func main() {
    router := httprouter.New()
    router.GET("/", Index)
    log.Fatal(http.ListenAndServe(":8080", router))
}
```

## TODO

- Recover stack trace for panic errors
- Add a easy way to initialize standard errors with an existing error. For example:

    ``` go
    //instead of
    err = errors.NewFromError(StatusInternalServer, err)
    // we could provide a function like this
    func InternalServerFromError(e error, setters ...paramsSetters)
    // and this can be used like this
    err := errors.InternalServerFromError(err)
    ```

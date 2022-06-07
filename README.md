[![codecov](https://codecov.io/gh/gofika/graceful/branch/main/graph/badge.svg)](https://codecov.io/gh/gofika/graceful)
[![Build Status](https://github.com/gofika/graceful/workflows/build/badge.svg)](https://github.com/gofika/graceful)
[![go.dev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/gofika/graceful)
[![Go Report Card](https://goreportcard.com/badge/github.com/gofika/graceful)](https://goreportcard.com/report/github.com/gofika/graceful)
[![Licenses](https://img.shields.io/github/license/gofika/graceful)](LICENSE)
[![donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.buymeacoffee.com/illi)

# Graceful

Graceful helpers for golang

## Basic Usage

### Installation

To get the package, execute:

```bash
go get github.com/gofika/graceful
```

To import this package, add the following line to your code:

```js
import "github.com/gofika/graceful";
```

### Example

Here is example usage.

```go
package main

import (
    "log"
    "net/http"

    "github.com/gofika/graceful"
)

func main() {
    ctx := context.Background()
    shutdown := graceful.NewShutdown(ctx)
  	r := http.NewServeMux()
  	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
      	w.Header().Set("Content-Type", "application/text")
      	w.Write([]byte("Success"))
  	}))
  	// Create a HTTP server and bind the router to it, and set wanted address
  	srv := &http.Server{
  		  Handler:      r,
  		  Addr:         ":8080",
  	}
    // Append closer for graceful shutdown
    shutdown.AppendGracefulClose(func() { srv.Close() })
    // Run graceful shutdown service
    go shutdown.Serve()
    // Run HTTP server. Server will graceful close util Ctrl+C signal
    if err := srv.ListenAndServe(); err != nil {
        if err != http.ErrServerClosed {
            log.Fatalf("server error: %s", err.Error())
        }
    }
}
```

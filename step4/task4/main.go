package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type params struct {
	Msg string
}

func parse(r *http.Request) params {
	msg := r.URL.Query().Get("msg")

	if msg == "" {
		return params{Msg: "empty"}
	}

	return params{Msg: msg}
}

func echo(w http.ResponseWriter, r *http.Request) {
	params := parse(r)
	fmt.Fprintf(w, "# %s", params.Msg)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", echo)

	server := &http.Server{Addr: "localhost:8080", Handler: mux}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Microsecond)
		defer cancel()

		if err := server.Shutdown(ctx); err != http.ErrServerClosed {
			fmt.Printf("shutdown server error: %v\n", err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("starting server error: %v\n", err)
	}
}

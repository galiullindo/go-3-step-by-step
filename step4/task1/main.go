package main

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type params struct {
	Name string
}

func parse(r *http.Request) params {
	name := r.URL.Query().Get("name")
	if name == "" {
		return params{Name: "stranger"}
	}

	onlyEn, _ := regexp.Match("^[a-zA-Z]+$", []byte(name))
	if !onlyEn {
		return params{Name: "dirty hacker"}
	}

	return params{Name: name}
}

func greeting(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	params := parse(r)
	fmt.Fprintf(w, "hello %s", params.Name)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", greeting)
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

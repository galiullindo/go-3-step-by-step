package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type UsernameKey string

const usernameKey UsernameKey = "username"

func Authorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username == "" || password == "" {
			w.Header().Add("WWW-Authenticate", "Basic realm=\"Restricted\"")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), usernameKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func answerHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Welcome, %s!", username)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/answer/", Authorization(answerHandler))

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

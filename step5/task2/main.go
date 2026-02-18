package main

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type Key string

const nameKey Key = "name"

func SetDefaultName(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "stranger"
		}

		ctx := context.WithValue(r.Context(), nameKey, name)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func Sanitize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, ok := r.Context().Value(nameKey).(string)
		if !ok {
			http.Error(w, "name lost", http.StatusInternalServerError)
			return
		}

		onlyEnglishLetters, _ := regexp.MatchString("^[a-zA-Z]+$", name)
		if !onlyEnglishLetters {
			name = "dirty hacker"
		}

		ctx := context.WithValue(r.Context(), nameKey, name)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name, ok := r.Context().Value(nameKey).(string)
	if !ok {
		http.Error(w, "name lost", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "hello %s", name)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", SetDefaultName(Sanitize(HelloHandler)))

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

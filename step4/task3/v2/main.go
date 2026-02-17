package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type SafeFibonacci struct {
	a       int
	b       int
	counter int
	mu      sync.Mutex
}

func NewSafeFibonacci() *SafeFibonacci {
	return &SafeFibonacci{0, 1, 0, sync.Mutex{}}
}

func (f *SafeFibonacci) Counter() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.counter
}

func (f *SafeFibonacci) Next() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.counter == 0 {
		f.counter++
		return f.a
	}
	if f.counter == 1 {
		f.counter++
		return f.b
	}

	f.counter++
	f.a, f.b = f.b, f.a+f.b
	return f.b
}

type FibonacciHandler struct {
	f *SafeFibonacci
}

func NewFibonacciHandler(f *SafeFibonacci) *FibonacciHandler {
	return &FibonacciHandler{f: f}
}

func (f *FibonacciHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fib := f.f.Next()
	fmt.Fprintf(w, "# %d", fib)
}

func metrics(f *SafeFibonacci) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fetches := f.Counter()
		fmt.Fprintf(w, "# %d", fetches)
	}
}

func main() {
	fibonacci := NewSafeFibonacci()
	handler := NewFibonacciHandler(fibonacci)

	mux := http.NewServeMux()
	mux.Handle("/", handler)
	mux.HandleFunc("/metrics", metrics(fibonacci))

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

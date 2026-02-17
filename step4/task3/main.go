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
	f       *SafeFibonacci
	counter int
	mu      sync.Mutex
}

func NewFibonacciHandler(f *SafeFibonacci) *FibonacciHandler {
	return &FibonacciHandler{f: f}
}

func (h *FibonacciHandler) Counter() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.counter
}

func (h *FibonacciHandler) Increment() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.counter++
}

func (h *FibonacciHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	h.Increment()

	fib := h.f.Next()
	fmt.Fprintf(w, "# %d", fib)
}

func metrics(h *FibonacciHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "# %d", h.Counter())
	}
}

func main() {
	fibonacci := NewSafeFibonacci()
	handler := NewFibonacciHandler(fibonacci)

	mux := http.NewServeMux()
	mux.Handle("/", handler)
	mux.HandleFunc("/metrics", metrics(handler))

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

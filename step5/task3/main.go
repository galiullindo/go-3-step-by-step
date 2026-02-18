package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type SafeFibonacci struct {
	a  int
	b  int
	n  int
	mu sync.Mutex
}

func NewSafeFibonacci() *SafeFibonacci {
	return &SafeFibonacci{0, 1, 0, sync.Mutex{}}
}

func (f *SafeFibonacci) Next() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.n == 0 {
		f.n++
		return 0
	}
	if f.n == 1 {
		f.n++
		return 1
	}

	f.n++
	f.a, f.b = f.b, f.a+f.b
	return f.b
}

type FibonacciHandler struct {
	f *SafeFibonacci
}

func NewFibonacciHandler(f *SafeFibonacci) *FibonacciHandler {
	return &FibonacciHandler{f: f}
}

func (h *FibonacciHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", h.f.Next())
}

type SafeCounter struct {
	value int
	mu    sync.Mutex
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{0, sync.Mutex{}}
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

var fibonacciCounter = NewSafeCounter()

func Metrics(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fibonacciCounter.Increment()
		next.ServeHTTP(w, r)
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "rpc_duration_milliseconds_count %d", fibonacciCounter.Value())
}

func main() {
	mux := http.NewServeMux()

	fibonacci := NewSafeFibonacci()
	fibonacciHandler := NewFibonacciHandler(fibonacci)
	mux.HandleFunc("/", Metrics(fibonacciHandler.ServeHTTP))

	mux.HandleFunc("/metrics", metricsHandler)

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

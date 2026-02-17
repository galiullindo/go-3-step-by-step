package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type safeFibonacci struct {
	a       int
	b       int
	counter int
	mu      sync.Mutex
}

func newSafeFibonacci() *safeFibonacci {
	return &safeFibonacci{0, 1, 0, sync.Mutex{}}
}

func (f *safeFibonacci) next() int {
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

type fibonacciHandler struct {
	f *safeFibonacci
}

func newFibonacciHandler(f *safeFibonacci) *fibonacciHandler {
	return &fibonacciHandler{f: f}
}

func (h *fibonacciHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fib := h.f.next()
	fmt.Fprintf(w, "# %d", fib)
}

func main() {
	fibonacci := newSafeFibonacci()
	handler := newFibonacciHandler(fibonacci)

	mux := http.NewServeMux()
	mux.Handle("/", handler)

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

package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"slices"
	"sort"
	"sync"
	"testing"
)

func TestSafeFibonacciСoncurrently(t *testing.T) {
	var tests = []struct {
		name     string
		times    int
		expected []int
	}{
		{
			name:     "next 0 numbers",
			times:    0,
			expected: []int(nil),
		},
		{
			name:     "next 1 numbers",
			times:    1,
			expected: []int{0},
		},
		{
			name:     "next 2 numbers",
			times:    2,
			expected: []int{0, 1},
		},
		{
			name:     "next 5 numbers",
			times:    5,
			expected: []int{0, 1, 1, 2, 3},
		},
		{
			name:     "next 10 numbers",
			times:    10,
			expected: []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			fibonacci := NewSafeFibonacci()

			channelFibonacciNumbers := make(chan int, test.times)
			for range test.times {
				go func() { channelFibonacciNumbers <- fibonacci.Next() }()
			}

			fibonacciNumbers := make([]int, test.times)
			for i := range test.times {
				fibonacciNumbers[i] = <-channelFibonacciNumbers
			}

			sort.Ints(fibonacciNumbers)
			if !slices.Equal(fibonacciNumbers, test.expected) {
				t.Errorf("unexpected fibonacci numbers: %v expected %v\n", fibonacciNumbers, test.expected)
			}
		})
	}
}

func TestFibonacciHandlerСoncurrently(t *testing.T) {
	var tests = []struct {
		name          string
		times         int
		expectedBodys []string
	}{
		{
			name:          "next 0 numbers",
			times:         0,
			expectedBodys: []string(nil),
		},
		{
			name:          "next 1 numbers",
			times:         1,
			expectedBodys: []string{"0"},
		},
		{
			name:          "next 2 numbers",
			times:         2,
			expectedBodys: []string{"0", "1"},
		},
		{
			name:          "next 5 numbers",
			times:         5,
			expectedBodys: []string{"0", "1", "1", "2", "3"},
		},
		{
			name:          "next 10 numbers",
			times:         10,
			expectedBodys: []string{"0", "1", "1", "2", "3", "5", "8", "13", "21", "34"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			fibonacci := NewSafeFibonacci()
			fibonacciHandler := NewFibonacciHandler(fibonacci)

			channelFibonacciBodys := make(chan string, test.times)
			for range test.times {
				go func(t *testing.T) {
					t.Helper()

					w := httptest.NewRecorder()
					fibonacciHandler.ServeHTTP(w, nil)

					response := w.Result()
					defer response.Body.Close()

					body, err := io.ReadAll(response.Body)
					if err != nil {
						t.Errorf("read response body: %v\n", err)
					}

					channelFibonacciBodys <- string(body)
				}(t)
			}

			fibonacciBodys := make([]string, test.times)
			for i := range test.times {
				fibonacciBodys[i] = <-channelFibonacciBodys
			}

			sort.Strings(test.expectedBodys)
			sort.Strings(fibonacciBodys)

			if !slices.Equal(fibonacciBodys, test.expectedBodys) {
				t.Errorf("unexpected fibonacci numbers: %v expected %v\n", fibonacciBodys, test.expectedBodys)
			}
		})
	}
}

func TestSafeCounterСoncurrently(t *testing.T) {
	var tests = []struct {
		name     string
		times    int
		expected int
	}{
		{name: "increment 0 times", times: 0, expected: 0},
		{name: "increment 10 times", times: 10, expected: 10},
		{name: "increment 100 times", times: 100, expected: 100},
		{name: "increment 1000 times", times: 1000, expected: 1000},
		{name: "increment 10000 times", times: 10000, expected: 10000},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			counter := NewSafeCounter()

			wg := &sync.WaitGroup{}
			for range test.times {
				wg.Go(func() {
					counter.Increment()
					_ = counter.Value()
				})
			}
			wg.Wait()

			if counter.Value() != test.expected {
				t.Errorf("unexpected counter value: %d expected %d\n", counter.Value(), test.expected)
			}
		})
	}
}

func TestMetricsСoncurrently(t *testing.T) {
	var tests = []struct {
		name     string
		times    int
		expected int
	}{
		{name: "0 requests", times: 0, expected: 0},
		{name: "10 requests", times: 10, expected: 10},
		{name: "100 requests", times: 100, expected: 100},
		{name: "1000 requests", times: 1000, expected: 1000},
		{name: "10000 requests", times: 10000, expected: 10000},
	}

	handlerWithMetrics := Metrics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fibonacciCounter = NewSafeCounter() // устанавливаем новый safe counter для каждого теста

			wg := &sync.WaitGroup{}
			for range test.times {
				wg.Go(func() {
					handlerWithMetrics.ServeHTTP(nil, nil)
				})
			}
			wg.Wait()

			if fibonacciCounter.Value() != test.expected {
				t.Errorf("unexpected metrics counter value: %d expected %d\n", fibonacciCounter.Value(), test.expected)
			}
		})
	}
}

func TestMetricsHandlerСoncurrently(t *testing.T) {
	var tests = []struct {
		name         string
		times        int
		expectedBody string
	}{
		{name: "0 requests to handler with metrics", times: 0, expectedBody: "rpc_duration_milliseconds_count 0"},
		{name: "10 requests to handler with metrics", times: 10, expectedBody: "rpc_duration_milliseconds_count 10"},
		{name: "100 requests to handler with metrics", times: 100, expectedBody: "rpc_duration_milliseconds_count 100"},
		{name: "1000 requests to handler with metrics", times: 1000, expectedBody: "rpc_duration_milliseconds_count 1000"},
		{name: "10000 requests to handler with metrics", times: 10000, expectedBody: "rpc_duration_milliseconds_count 10000"},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {})
	handlerWithMetrics := Metrics(nextHandler)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fibonacciCounter = NewSafeCounter()

			wg := &sync.WaitGroup{}
			for range test.times {
				wg.Go(func() {
					handlerWithMetrics.ServeHTTP(nil, nil)
				})
			}
			wg.Wait()

			w := httptest.NewRecorder()
			metricsHandler(w, nil)

			response := w.Result()
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Printf("read response body: %v\n", err)
				return
			}

			if string(body) != test.expectedBody {
				t.Errorf("unexpected metrics body: %q expected %q\n", string(body), test.expectedBody)
			}
		})
	}
}

func TestServer(t *testing.T) {
	type step struct {
		method     string
		handlerUrl string
		expected   string
	}

	var tests = []struct {
		name  string
		steps []step
	}{
		{
			name: "",
			steps: []step{
				{method: http.MethodGet, handlerUrl: "/metrics", expected: "rpc_duration_milliseconds_count 0"},
				{method: http.MethodGet, handlerUrl: "/metrics", expected: "rpc_duration_milliseconds_count 0"},
				{method: http.MethodGet, handlerUrl: "/", expected: "0"},
				{method: http.MethodGet, handlerUrl: "/", expected: "1"},
				{method: http.MethodGet, handlerUrl: "/", expected: "1"},
				{method: http.MethodGet, handlerUrl: "/", expected: "2"},
				{method: http.MethodGet, handlerUrl: "/", expected: "3"},
				{method: http.MethodGet, handlerUrl: "/", expected: "5"},
				{method: http.MethodGet, handlerUrl: "/", expected: "8"},
				{method: http.MethodGet, handlerUrl: "/metrics", expected: "rpc_duration_milliseconds_count 7"},
				{method: http.MethodGet, handlerUrl: "/metrics", expected: "rpc_duration_milliseconds_count 7"},
				{method: http.MethodGet, handlerUrl: "/", expected: "13"},
				{method: http.MethodGet, handlerUrl: "/", expected: "21"},
				{method: http.MethodGet, handlerUrl: "/", expected: "34"},
				{method: http.MethodGet, handlerUrl: "/metrics", expected: "rpc_duration_milliseconds_count 10"},
				{method: http.MethodGet, handlerUrl: "/metrics", expected: "rpc_duration_milliseconds_count 10"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fibonacciCounter = NewSafeCounter() // устанавливаем новый safe counter для каждого теста

			mux := http.NewServeMux()

			fibonacci := NewSafeFibonacci()
			fibonacciHandler := NewFibonacciHandler(fibonacci)

			mux.HandleFunc("/", Metrics(fibonacciHandler.ServeHTTP))
			mux.HandleFunc("/metrics", metricsHandler)

			server := httptest.NewServer(mux)
			defer server.Close()

			for _, step := range test.steps {
				request, err := http.NewRequest(step.method, server.URL+step.handlerUrl, nil)
				if err != nil {
					t.Fatalf("new request: %v\n", err)
				}

				response, err := http.DefaultClient.Do(request)
				if err != nil {
					t.Fatalf("response: %v\n", err)
				}
				defer response.Body.Close()

				body, err := io.ReadAll(response.Body)
				if err != nil {
					t.Fatalf("read respone body: %v\n", err)
				}

				if string(body) != step.expected {
					t.Errorf("unexpected body: %q expected %q\n", string(body), step.expected)
				}
			}
		})
	}
}

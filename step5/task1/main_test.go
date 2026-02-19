package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	var tests = []struct {
		name        string
		w           *httptest.ResponseRecorder
		r           *http.Request
		expectedLog string
	}{
		{
			name:        "Method GET",
			w:           httptest.NewRecorder(),
			r:           httptest.NewRequest(http.MethodGet, "/hello", nil),
			expectedLog: "level=INFO msg=\"incoming request\" method=GET path=/hello\n",
		},
		{
			name:        "Method POST",
			w:           httptest.NewRecorder(),
			r:           httptest.NewRequest(http.MethodPost, "/hello", nil),
			expectedLog: "level=INFO msg=\"incoming request\" method=POST path=/hello\n",
		},
	}

	log := bytes.NewBuffer(nil)
	logHandler := slog.NewTextHandler(log, nil)
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	handlerWithLogging := Logger(nextHandler)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer log.Truncate(0)

			handlerWithLogging.ServeHTTP(test.w, test.r)
			if !strings.Contains(log.String(), test.expectedLog) {
				t.Errorf("unexpected log: %q expected %q\n", log.String(), test.expectedLog)
			}
		})
	}
}

func TestHelloHandler(t *testing.T) {
	var tests = []struct {
		name               string
		w                  *httptest.ResponseRecorder
		r                  *http.Request
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "Method GET",
			w:                  httptest.NewRecorder(),
			r:                  httptest.NewRequest(http.MethodGet, "/hello", nil),
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Hello, middleware!",
		},
		{
			name:               "Method POST",
			w:                  httptest.NewRecorder(),
			r:                  httptest.NewRequest(http.MethodPost, "/hello", nil),
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Hello, middleware!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			helloHandler(test.w, test.r)
			response := test.w.Result()
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("read respone body: %v\n", err)
			}

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("unexpected status code: %q expected %q\n", response.Status, test.expectedStatusCode)
			}

			if string(body) != test.expectedBody {
				t.Errorf("unexpected body: %q expected %q\n", string(body), test.expectedBody)
			}
		})
	}
}

func TestServer(t *testing.T) {
	var tests = []struct {
		name               string
		method             string
		handlerUrl         string
		expectedBody       string
		expectedStatusCode int
		expectedLog        string
	}{
		{
			name:               "Method GET",
			method:             http.MethodGet,
			handlerUrl:         "/hello",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Hello, middleware!",
			expectedLog:        "level=INFO msg=\"incoming request\" method=GET path=/hello\n",
		},
		{
			name:               "Method POST",
			method:             http.MethodPost,
			handlerUrl:         "/hello",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Hello, middleware!",
			expectedLog:        "level=INFO msg=\"incoming request\" method=POST path=/hello\n",
		},
	}

	log := bytes.NewBuffer(nil)
	logHandler := slog.NewTextHandler(log, nil)
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	handler := Logger(mux)

	server := httptest.NewServer(handler)
	defer server.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log.Truncate(0)

			request, err := http.NewRequest(test.method, server.URL+test.handlerUrl, nil)
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

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("unexpected status code: %q expected %q\n", response.Status, test.expectedStatusCode)
			}

			if string(body) != test.expectedBody {
				t.Errorf("unexpected body: %q expected %q\n", string(body), test.expectedBody)
			}

			if !strings.Contains(log.String(), test.expectedLog) {
				t.Errorf("unexpected log: %q expected %q\n", log.String(), test.expectedLog)
			}
		})
	}
}

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func nextHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Helper()

		name, ok := r.Context().Value(nameKey).(string)
		if !ok {
			t.Errorf("lost name")
		}

		fmt.Fprintf(w, "%s", name)
	}
}

func TestSetDefaultName(t *testing.T) {
	var tests = []struct {
		name         string
		w            *httptest.ResponseRecorder
		r            *http.Request
		expectedBody string
	}{
		{
			name:         "hello",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello", nil),
			expectedBody: "stranger",
		},
		{
			name:         "hello?",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?", nil),
			expectedBody: "stranger",
		},
		{
			name:         "hello?name=",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=", nil),
			expectedBody: "stranger",
		},
		{
			name:         "hello?name=abc",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=abc", nil),
			expectedBody: "abc",
		},
		{
			name:         "hello?name=ABC",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=ABC", nil),
			expectedBody: "ABC",
		},
		{
			name:         "hello?name=123",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=123", nil),
			expectedBody: "123",
		},
		{
			name:         "hello?name=абв",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=абв", nil),
			expectedBody: "абв",
		},
		{
			name:         "hello?name=АБВ",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=АБВ", nil),
			expectedBody: "АБВ",
		},
		{
			name:         "hello?name=abc123",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=abc123", nil),
			expectedBody: "abc123",
		},
		{
			name:         "hello?name=abcабв",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=abcабв", nil),
			expectedBody: "abcабв",
		},
		{
			name:         "hello?name=abcАБВ",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=abcАБВ", nil),
			expectedBody: "abcАБВ",
		},
	}

	handler := SetDefaultName(nextHandler(t))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			handler.ServeHTTP(test.w, test.r)
			response := test.w.Result()
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("read respone body: %v\n", err)
			}

			if string(body) != test.expectedBody {
				t.Errorf("unexpected name: %q expected %q\n", string(body), test.expectedBody)
			}
		})
	}
}

func TestSanitize(t *testing.T) {
	var tests = []struct {
		name         string
		w            *httptest.ResponseRecorder
		r            *http.Request
		ctx          context.Context
		expectedBody string
	}{
		{
			name:         "internal server error",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello", nil),
			ctx:          context.Background(),
			expectedBody: "name lost\n",
		},
		{
			name:         "name stranger",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=", nil),
			ctx:          context.WithValue(context.Background(), nameKey, "stranger"),
			expectedBody: "stranger",
		},
		{
			name:         "hello?name=abc",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=abc", nil),
			ctx:          context.WithValue(context.Background(), nameKey, "abc"),
			expectedBody: "abc",
		},
		{
			name:         "hello?name=ABC",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=ABC", nil),
			ctx:          context.WithValue(context.Background(), nameKey, "ABC"),
			expectedBody: "ABC",
		},
		{
			name:         "hello?name=123",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=123", nil),
			ctx:          context.WithValue(context.Background(), nameKey, "123"),
			expectedBody: "dirty hacker",
		},
		{
			name:         "hello?name=абв",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=абв", nil),
			ctx:          context.WithValue(context.Background(), nameKey, "абв"),
			expectedBody: "dirty hacker",
		},
		{
			name:         "hello?name=АБВ",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=АБВ", nil),
			ctx:          context.WithValue(context.Background(), nameKey, "АБВ"),
			expectedBody: "dirty hacker",
		},
		{
			name:         "hello?name=abc123",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=abc123", nil),
			ctx:          context.WithValue(context.Background(), nameKey, "abc123"),
			expectedBody: "dirty hacker",
		},
		{
			name:         "hello?name=abcабв",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=abcабв", nil),
			ctx:          context.WithValue(context.Background(), nameKey, "abcабв"),
			expectedBody: "dirty hacker",
		},
		{
			name:         "hello?name=abcАБВ",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=abcАБВ", nil),
			ctx:          context.WithValue(context.Background(), nameKey, "abcАБВ"),
			expectedBody: "dirty hacker",
		},
	}

	handler := Sanitize(nextHandler(t))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			handler.ServeHTTP(test.w, test.r.WithContext(test.ctx))
			response := test.w.Result()
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("read respone body: %v\n", err)
			}

			if string(body) != test.expectedBody {
				t.Errorf("unexpected name: %q expected %q\n", string(body), test.expectedBody)
			}
		})
	}
}

func TestHelloHandler(t *testing.T) {
	var tests = []struct {
		name         string
		w            *httptest.ResponseRecorder
		r            *http.Request
		ctx          context.Context
		expectedBody string
	}{
		{
			name:         "internal server error",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello", nil),
			ctx:          context.Background(),
			expectedBody: "name lost\n",
		},
		{
			name:         "name stranger",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=", nil),
			ctx:          context.WithValue(context.Background(), nameKey, "stranger"),
			expectedBody: "hello stranger",
		},
		{
			name:         "name dirty hacker",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=abcABCабв123АБВ", nil),
			ctx:          context.WithValue(context.Background(), nameKey, "dirty hacker"),
			expectedBody: "hello dirty hacker",
		},
		{
			name:         "hello?name=abc",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=abc", nil),
			ctx:          context.WithValue(context.Background(), nameKey, "abc"),
			expectedBody: "hello abc",
		},
		{
			name:         "hello?name=ABC",
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/hello?name=ABC", nil),
			ctx:          context.WithValue(context.Background(), nameKey, "ABC"),
			expectedBody: "hello ABC",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			HelloHandler(test.w, test.r.WithContext(test.ctx))
			response := test.w.Result()
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("read respone body: %v\n", err)
			}

			if string(body) != test.expectedBody {
				t.Errorf("unexpected name: %q expected %q\n", string(body), test.expectedBody)
			}
		})
	}
}

func TestServer(t *testing.T) {
	var tests = []struct {
		name         string
		method       string
		handlerUrl   string
		r            *http.Request
		expectedBody string
	}{
		{
			name:         "hello",
			method:       http.MethodGet,
			handlerUrl:   "/hello",
			expectedBody: "hello stranger",
		},
		{
			name:         "hello?",
			method:       http.MethodGet,
			handlerUrl:   "/hello?",
			expectedBody: "hello stranger",
		},
		{
			name:         "hello?name=",
			method:       http.MethodGet,
			handlerUrl:   "/hello?name=",
			expectedBody: "hello stranger",
		},
		{
			name:         "hello?name=abc",
			method:       http.MethodGet,
			handlerUrl:   "/hello?name=abc",
			expectedBody: "hello abc",
		},
		{
			name:         "hello?name=ABC",
			method:       http.MethodGet,
			handlerUrl:   "/hello?name=ABC",
			expectedBody: "hello ABC",
		},
		{
			name:         "hello?name=123",
			method:       http.MethodGet,
			handlerUrl:   "/hello?name=123",
			expectedBody: "hello dirty hacker",
		},
		{
			name:         "hello?name=абв",
			method:       http.MethodGet,
			handlerUrl:   "/hello?name=абв",
			expectedBody: "hello dirty hacker",
		},
		{
			name:         "hello?name=АБВ",
			method:       http.MethodGet,
			handlerUrl:   "/hello?name=АБВ",
			expectedBody: "hello dirty hacker",
		},
		{
			name:         "hello?name=abc123",
			method:       http.MethodGet,
			handlerUrl:   "/hello?name=abc123",
			expectedBody: "hello dirty hacker",
		},
		{
			name:         "hello?name=abcабв",
			method:       http.MethodGet,
			handlerUrl:   "/hello?name=abcабв",
			expectedBody: "hello dirty hacker",
		},
		{
			name:         "hello?name=abcАБВ",
			method:       http.MethodGet,
			handlerUrl:   "/hello?name=abcАБВ",
			expectedBody: "hello dirty hacker",
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", SetDefaultName(Sanitize(HelloHandler)))

	server := httptest.NewServer(mux)
	defer server.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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

			if string(body) != test.expectedBody {
				t.Errorf("unexpected body: %q expected %q\n", string(body), test.expectedBody)
			}
		})
	}
}

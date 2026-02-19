package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthorization(t *testing.T) {
	type header struct {
		name  string
		value string
	}

	var tests = []struct {
		name               string
		username           string
		password           string
		expectedStatusCode int
		expectedBody       string
		expectedHeaders    []header
		expecteUsername    string
	}{
		{
			name:               "without username and password",
			username:           "",
			password:           "",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "Unauthorized\n",
			expectedHeaders:    []header{{"WWW-Authenticate", "Basic realm=\"Restricted\""}},
		},
		{
			name:               "without username and password",
			username:           "",
			password:           "123",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "Unauthorized\n",
			expectedHeaders:    []header{{"WWW-Authenticate", "Basic realm=\"Restricted\""}},
		},
		{
			name:               "without password",
			username:           "abc",
			password:           "",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "Unauthorized\n",
			expectedHeaders:    []header{{"WWW-Authenticate", "Basic realm=\"Restricted\""}},
		},
		{
			name:               "username and password abc:123",
			username:           "abc",
			password:           "123",
			expectedStatusCode: http.StatusOK,
			expecteUsername:    "abc",
		},
		{
			name:               "username and password ABCDEF:123456",
			username:           "ABCDEF",
			password:           "123456",
			expectedStatusCode: http.StatusOK,
			expecteUsername:    "ABCDEF",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				username, ok := r.Context().Value(usernameKey).(string)
				if !ok {
					t.Errorf("username not found in context")
				}
				if username != test.username {
					t.Errorf("unexpected username %q expected %q\n", username, test.username)
				}
				w.WriteHeader(http.StatusOK)
			})

			handlerWithAuthorization := Authorization(nextHandler)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/answer/", nil)
			r.SetBasicAuth(test.username, test.password)

			handlerWithAuthorization.ServeHTTP(w, r)

			response := w.Result()
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Printf("read response body: %v\n", err)
				return
			}

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("unexpected status code %d expected %d", response.StatusCode, test.expectedStatusCode)
			}

			if string(body) != test.expectedBody {
				t.Errorf("unexpected body %q expected %q\n", string(body), test.expectedBody)
			}

			for _, header := range test.expectedHeaders {
				value := response.Header.Get(header.name)
				if value != header.value {
					t.Errorf("unexpected header %q value %q expected %q\n", header.name, value, header.value)
				}
			}
		})
	}
}

func TestAnswerHandler(t *testing.T) {
	var tests = []struct {
		name               string
		ctx                context.Context
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "internal server error",
			ctx:                context.Background(),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "InternalServerError\n",
		},
		{
			name:               "username abc",
			ctx:                context.WithValue(context.Background(), usernameKey, "abc"),
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Welcome, abc!",
		},
		{
			name:               "username ABCDEF",
			ctx:                context.WithValue(context.Background(), usernameKey, "ABCDEF"),
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Welcome, ABCDEF!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/answer/", nil)

			answerHandler(w, r.WithContext(test.ctx))

			response := w.Result()
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Printf("read response body: %v\n", err)
				return
			}

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("unexpected status code %d expected %d", response.StatusCode, test.expectedStatusCode)
			}

			if string(body) != test.expectedBody {
				t.Errorf("unexpected body %q expected %q\n", string(body), test.expectedBody)
			}
		})
	}
}

func TestServer(t *testing.T) {
	type header struct {
		name  string
		value string
	}

	var tests = []struct {
		name               string
		method             string
		handlerUrl         string
		username           string
		password           string
		expectedStatusCode int
		expectedBody       string
		expectedHeaders    []header
	}{
		{
			name:               "without username and password",
			method:             http.MethodGet,
			handlerUrl:         "/answer/",
			username:           "",
			password:           "",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "Unauthorized\n",
			expectedHeaders:    []header{{"WWW-Authenticate", "Basic realm=\"Restricted\""}},
		},
		{
			name:               "without username and password",
			method:             http.MethodGet,
			handlerUrl:         "/answer/",
			username:           "",
			password:           "123",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "Unauthorized\n",
			expectedHeaders:    []header{{"WWW-Authenticate", "Basic realm=\"Restricted\""}},
		},
		{
			name:               "without password",
			method:             http.MethodGet,
			handlerUrl:         "/answer/",
			username:           "abc",
			password:           "",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "Unauthorized\n",
			expectedHeaders:    []header{{"WWW-Authenticate", "Basic realm=\"Restricted\""}},
		},
		{
			name:               "username and password abc:123",
			method:             http.MethodGet,
			handlerUrl:         "/answer/",
			username:           "abc",
			password:           "123",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Welcome, abc!",
		},
		{
			name:               "username and password ABCDEF:123456",
			method:             http.MethodGet,
			handlerUrl:         "/answer/",
			username:           "ABCDEF",
			password:           "123456",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Welcome, ABCDEF!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/answer/", Authorization(answerHandler))

			server := httptest.NewServer(mux)

			request, err := http.NewRequest(test.method, server.URL+test.handlerUrl, nil)
			if err != nil {
				t.Fatalf("new request: %v\n", err)
			}
			request.SetBasicAuth(test.username, test.password)

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatalf("response: %v\n", err)
			}
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Printf("read response body: %v\n", err)
				return
			}

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("unexpected status code %d expected %d", response.StatusCode, test.expectedStatusCode)
			}

			if string(body) != test.expectedBody {
				t.Errorf("unexpected body %q expected %q\n", string(body), test.expectedBody)
			}

			for _, header := range test.expectedHeaders {
				value := response.Header.Get(header.name)
				if value != header.value {
					t.Errorf("unexpected header %q value %q expected %q\n", header.name, value, header.value)
				}
			}
		})
	}
}

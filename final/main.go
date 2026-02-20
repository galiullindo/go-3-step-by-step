package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Store struct {
	nextUserID int
	users      map[int]User
	mu         sync.Mutex
}

func NewStore() *Store {
	return &Store{nextUserID: 1, users: map[int]User{}, mu: sync.Mutex{}}
}

func (s *Store) CreateUser(name string, age int) User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := User{ID: s.nextUserID, Name: name, Age: age}
	s.users[s.nextUserID] = user
	s.nextUserID++

	return user
}

func (s *Store) GetUser(id int) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[id]

	return user, ok
}

type responseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrappedW := &responseWriter{ResponseWriter: w, StatusCode: 200}
		start := time.Now()

		next.ServeHTTP(wrappedW, r)

		duration := time.Since(start)
		logger.Info(
			"http request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", wrappedW.StatusCode),
			slog.Duration("duration", duration),
		)
	})
}

func createUserHandler(storage *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		storage.CreateUser(user.Name, user.Age)
	})
}

func getUserHandler(storage *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strId := r.PathValue("id")
		if strId == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(strId)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		user, found := storage.GetUser(id)
		if !found {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	})
}

func main() {
	mux := http.NewServeMux()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	storage := NewStore()

	mux.Handle("POST /users", createUserHandler(storage))
	mux.Handle("GET /users/{id}", getUserHandler(storage))

	handler := loggingMiddleware(logger, mux)
	http.ListenAndServe("localhost:8080", handler)
}

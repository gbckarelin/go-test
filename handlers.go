package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	_ "yanego/docs"
)

// @Summary Create a new task
// @Description Create a new task for the authenticated user
// @Tags tasks
// @Produce  json
// @Success 201 {string} string "Task ID"
// @Failure 500 {string} string "Internal server error"
// @Router /task [post]
func handleTask(w http.ResponseWriter, r *http.Request) {

	username := "loh"

	// username, auth := r.Context().Value("username").(string)
	// if !auth || username == "" {
	// 	http.Error(w, "Unnauthorized", http.StatusUnauthorized)
	// 	return
	// }

	taskID := store.CreateTask(username)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"task_id": taskID})
}

// @Summary Get task status
// @Description Get the status of a specific task by its ID
// @Tags tasks
// @Produce  json
// @Param   taskID  path   string  true  "Task ID"
// @Success 200 {object} TaskStatus "Task status"
// @Failure 404 {string} string "Task not found"
// @Router /status/{taskID} [get]
func handleStatus(w http.ResponseWriter, r *http.Request) {

	taskID := chi.URLParam(r, "taskID")
	status, err := store.GetTaskStatus(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// @Summary Get task result
// @Description Get the result of a specific task by its ID
// @Tags tasks
// @Produce  json
// @Param   taskID  path   string  true  "Task ID"
// @Success 200 {string} string "Task result"
// @Failure 404 {string} string "Task not found"
// @Failure 202 {string} string "Task not ready"
// @Router /result/{taskID} [get]
func handleRequest(w http.ResponseWriter, r *http.Request) {

	taskID := chi.URLParam(r, "taskID")
	result, err := store.GetTaskResult(taskID)

	if err != nil {
		if err.Error() == "task not ready" {
			http.Error(w, err.Error(), http.StatusAccepted)
		} else {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

// @Summary Register a new user
// @Description Create a new user with a username and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user  body   User  true  "User data"
// @Success 201 {string} string "User registered successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 409 {string} string "User already exists"
// @Router /register [post]
func handleRegister(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := store.RegisterUser(user.Username, user.Password); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

// @Summary Login a user
// @Description Authenticate user and return a JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user  body   User  true  "User data"
// @Success 200 {string} string "JWT token"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Invalid username or password"
// @Router /login [post]
func handleLogin(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	token, err := store.AuthenticateUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}

func auth(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHead := r.Header.Get("Authorization")
		if authHead == "" {
			http.Error(w, "Unnauthorized", http.StatusUnauthorized)
			return
		}

		splits := strings.Split(authHead, " ")
		if splits[0] != "Bearer" || len(splits) != 2 {
			http.Error(w, "Unnauthorized", http.StatusUnauthorized)
			return
		}

		token := splits[1]
		username, valid := store.ValidateToken(token)
		if !valid {
			http.Error(w, "Unnauthorized", http.StatusUnauthorized)
			return
		}

		part := r.Context()
		part = context.WithValue(part, "username", username)
		n.ServeHTTP(w, r.WithContext(part))
	})
}

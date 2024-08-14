package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func handleTask(w http.ResponseWriter, r *http.Request) {

	username := "loh"

	taskID := store.CreateTask(username)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"task_id": taskID})
}

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

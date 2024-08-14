package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type TaskStatus struct {
	Status string `json:"status"`
	Result string `json:"result,omitempty"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func handleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	taskID := uuid.New().String()
	tasks[taskID] = &TaskStatus{
		Status: "in_progress",
	}
	go func() {
		time.Sleep(5 * time.Second)
		tasks[taskID].Status = "ready"
		tasks[taskID].Result = "mysor"
	}()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"task_id": taskID})
	fmt.Println("mysor")
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Path[len("/status/"):]
	task, exists := tasks[taskID]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": task.Status})
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Path[len("/result/"):]
	task, exists := tasks[taskID]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	if task.Status != "ready" {
		http.Error(w, "Task not ready", http.StatusAccepted)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result": task.Result})
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Nill Username or password", http.StatusBadRequest)
		return
	}

	if _, exists := users[user.Username]; exists {
		http.Error(w, "User already registred", http.StatusConflict)
		return
	}

	users[user.Username] = &user
	w.WriteHeader(http.StatusCreated)

}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var login User
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	storedUser, exists := users[login.Username]
	if !exists || storedUser.Password != login.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token := uuid.New().String()
	tokens[token] = login.Username
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}

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

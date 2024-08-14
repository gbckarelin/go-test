package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Storage struct {
	users  map[string]*User
	tasks  map[string]*TaskStatus
	tokens map[string]string
}

type TaskStatus struct {
	Status string `json:"status"`
	Result string `json:"result,omitempty"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewStorage() *Storage {
	return &Storage{
		users:  make(map[string]*User),
		tasks:  make(map[string]*TaskStatus),
		tokens: make(map[string]string),
	}
}

//регистрация пользователя
func (s *Storage) RegisterUser(username, password string) error {
	if _, exists := s.users[username]; exists {
		return fmt.Errorf("user already exists")
	}
	s.users[username] = &User{Username: username, Password: password}
	return nil
}

//аутентификация
func (s *Storage) AuthenticateUser(username, password string) (string, error) {
	user, exists := s.users[username]
	if !exists || user.Password != password {
		return "", fmt.Errorf("invalid username or password")
	}
	token := uuid.New().String()
	s.tokens[token] = username
	return token, nil
}

//получение Id таски
func (s *Storage) CreateTask(username string) string {
	taskID := uuid.New().String()
	s.tasks[taskID] = &TaskStatus{Status: "in_progress"}
	go func() {
		// эмуляция
		time.Sleep(5 * time.Second)
		s.tasks[taskID].Status = "ready"
		s.tasks[taskID].Result = "mysor"
	}()
	return taskID
}

// получение статуса таски
func (s *Storage) GetTaskStatus(taskID string) (*TaskStatus, error) {
	task, exists := s.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

//получение результата таски
func (s *Storage) GetTaskResult(taskID string) (string, error) {
	task, exists := s.tasks[taskID]
	if !exists {
		return "", fmt.Errorf("task not found")
	}
	if task.Status != "ready" {
		return "", fmt.Errorf("task not ready")
	}
	return task.Result, nil
}

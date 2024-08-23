package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Storage struct {
	users    map[string]*User
	sessions map[string]*Session
	tasks    map[string]*TaskStatus
	tokens   map[string]string
}

type Session struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
}

type TaskStatus struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewStorage() *Storage {
	return &Storage{
		users:    make(map[string]*User),
		sessions: make(map[string]*Session),
		tasks:    make(map[string]*TaskStatus),
		tokens:   make(map[string]string),
	}
}

// регистрация пользователя
func (s *Storage) RegisterUser(username, password string) error {
	if _, exists := s.users[username]; exists {
		return fmt.Errorf("user already exists")
	}
	id := uuid.New().String()
	s.users[username] = &User{ID: id, Username: username, Password: password}
	return nil
}

// аутентификация
func (s *Storage) AuthenticateUser(username, password string) (string, error) {
	user, exists := s.users[username]
	if !exists || user.Password != password {
		return "", fmt.Errorf("invalid username or password")
	}
	token := uuid.New().String()
	session := &Session{
		UserID:    user.ID,
		SessionID: token,
	}
	s.sessions[token] = session
	s.tokens[token] = username
	return token, nil
}

func (s *Storage) ValidateToken(token string) (string, bool) {
	username, exists := s.tokens[token]
	return username, exists
}

// получение Id таски
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

// получение результата таски
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

func (s *Storage) doTask(taskID string) {
	time.Sleep(5 * time.Second)
	s.tasks[taskID].Status = "ready"
	s.tasks[taskID].Result = "mysor"
}

package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var store = NewStorage()

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/register", handleRegister)
	r.Post("/login", handleLogin)
	r.Post("/task", handleTask)
	r.Get("/status/{taskID}", handleStatus)
	r.Get("/result/{taskID}", handleRequest)

	fmt.Println("Server start 8000")
	http.ListenAndServe(":8000", r)
}

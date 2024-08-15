package main

import (
	"fmt"
	"net/http"

	_ "yanego/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

var store = NewStorage()

// @title Swagger Example API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/register", handleRegister)
	r.Post("/login", handleLogin)
	r.Post("/task", handleTask)
	r.Get("/status/{taskID}", handleStatus)
	r.Get("/result/{taskID}", handleRequest)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	fmt.Println("Server start 8000")
	http.ListenAndServe(":8000", r)
}

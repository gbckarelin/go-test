package main

import (
	"fmt"
	"net/http"
)

var tasks = make(map[string]*TaskStatus)
var users = make(map[string]*User)
var tokens = make(map[string]string)

func main() {
	http.HandleFunc("/task", handleTask)
	http.HandleFunc("/status/", handleStatus)
	http.HandleFunc("/result/", handleRequest)
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/login", handleLogin)

	fmt.Println("Server start 8000")
	http.ListenAndServe(":8000", nil)
}

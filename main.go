package main

import (
	"fmt"
	"net/http"
)

var tasks = make(map[string]*TaskStatus)

func main() {
	http.HandleFunc("/task", handleTask)
	http.HandleFunc("/status/", handleStatus)
	http.HandleFunc("/result/", handleRequest)

	fmt.Println("Server start 8080")
	http.ListenAndServe(":8080", nil)
}

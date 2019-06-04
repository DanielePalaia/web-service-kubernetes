package main

import (
	"net/http"
	"time"
	"web-service-kubernetes/servicelog"
)

// Main entry point of the service listening on 8080
func main() {
	//fmt.Println("starting...")
	logger := servicelog.GetInstance()
	logger.Println(time.Now().UTC(), "Starting service")
	// routes defined in routes.go
	router := NewRouter()
	http.ListenAndServe(":8080", router)
}

package main

import (
	"log"
	"net/http"
)
func main(){
	const (
		maxWorkers = 4
		maxQueueSize = 20
		port = ":8081"
	)
	jobQueue := make(chan Job, maxQueueSize)
	dispatcher := NewDispatcher(jobQueue, maxWorkers)

	dispatcher.Run()
	http.HandleFunc("/fib", RequestHandler(jobQueue))
	log.Fatal(http.ListenAndServe(port, nil))
}
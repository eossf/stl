package main

import (
	"log"
	"net/http"
)

func main() {
	testMongodb()
	log.Println("Start HTTP server")
	// start http server
	router := NewRouter(AllRoutes())
	log.Fatal(http.ListenAndServe(":8080", router))
}

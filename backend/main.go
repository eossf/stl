package main

import (
	"log"
	"net/http"
)

func main() {
	// check the mongodb with track id = 1
	getTrack("1")
	log.Println("Start HTTP server")
	// start http server
	router := NewRouter(AllRoutes())
	log.Fatal(http.ListenAndServe(":8080", router))
}

package main

import (
	"log"
	"net/http"
)

func main() {
	// check the mongodb with track id = 1
	getTrack(1)
	log.Println("Start HTTP server")
	// start http server
	router := NewRouter(AllRoutes())

	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Credentials", "true")
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Headers", "Origin,Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,locale")
			header.Set("Access-Control-Allow-Methods", "GET, PUT, POST, OPTIONS")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	log.Fatal(http.ListenAndServe(": $os.Getegid('PORT_STL_BACKEND')", router))
}

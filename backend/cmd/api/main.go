package main

import (
	"log"
	"net/http"
	handler "testy/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/test", handler.SubmitHandler)

	log.Println("Server locked in at http://127.0.0.1:8090")
	if err := http.ListenAndServe("0.0.0.0:8090", handler.CORSMiddleware(mux)); err != nil {
		log.Fatalf("Server crashed: %v", err)
	}
}

package main

import (
	"context"
	"log"
	"testy/server"
)

func main() {
	server, err := server.NewServer(context.Background(), "0.0.0.0:8090", "postgres://dev:dev@postgres:5432/testy?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
		return
	}
	server.Serve()
}

package main

import (
	"log"

	"github.com/shakareem/review-assigner/internal/api"
	"github.com/shakareem/review-assigner/internal/server"
	"github.com/shakareem/review-assigner/internal/storage"
)

func main() {
	storage, err := storage.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	handler := api.NewHandler(storage)
	server := server.NewServer(handler)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/gorilla/mux"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	//
	// Mux router part
	//
	router := mux.NewRouter()

	s := Service{
		Message{getEnv("MSG", "Hello world!")},
	}

	// message = getEnv("MSG", "Hello world!")

	router.HandleFunc("/", s.showHello).Methods("GET")

	port := fmt.Sprintf(":%s", getEnv("PORT", "8080"))

	log.Printf("Listening on %s...", port)
	log.Fatal(http.ListenAndServe(port, router))
}

type Message struct {
	Message string `json:"message"`
}

type Service struct {
	message Message
}

func (s *Service) showHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&s.message)
}

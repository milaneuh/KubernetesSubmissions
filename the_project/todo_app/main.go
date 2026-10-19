package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/healthz", healthzHandler)

	log.Println("Server started in port " + port)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatal("Could not start the server\n", err)
	}
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}

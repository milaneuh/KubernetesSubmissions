package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

var counter int

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", PingHandler)

	log.Println("Server started in port " + port)
	http.ListenAndServe(":"+port, nil)
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	counter += 1
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong " + strconv.Itoa(counter)))
}

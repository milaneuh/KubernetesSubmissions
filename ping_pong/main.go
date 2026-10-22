package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var counter int

func main() {
	port := os.Getenv("PORT")
	path := filePath()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", PingHandler)

	log.Println("Server started in port " + port)
	http.ListenAndServe(":"+port, nil)
}

func filePath() string {
	if p := os.Getenv("PONG_FILE"); p != "" {
		return p
	}
	return "/usr/src/app/files/pingpong_output.txt"
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	counter += 1

	err := os.WriteFile(filePath(), []byte(strconv.Itoa(counter)), 0644)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong " + strconv.Itoa(counter)))
}

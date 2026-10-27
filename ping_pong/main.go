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

	// Full paths: the ingress forwards /pingpong... without stripping the prefix
	http.HandleFunc("/pingpong", PingHandler)
	http.HandleFunc("/pingpong/pings", PingsHandler)

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

// PingsHandler returns the current count without incrementing it.
func PingsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(counter)))
}

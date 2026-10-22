package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var hash string

func main() {
	port := os.Getenv("PORT")
	path := filePath()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Fatal(err)
	}

	err := GenerateRandomHash(512)
	if err != nil {
		log.Fatal(err)
	}

	if port != "" {
		http.HandleFunc("/", HashHandler)

		log.Println("Server started in port " + port)
		http.ListenAndServe(":"+port, nil)
	} else {
		PeriodicallyPrintHash(hash)
	}
}

func filePath() string {
	if p := os.Getenv("LOG_FILE"); p != "" {
		return p
	}
	return "/usr/src/app/files/log_output.txt"
}

func pongFilePath() string {
	if p := os.Getenv("PONG_FILE"); p != "" {
		return p
	}
	return "/usr/src/app/files/pingpong_output.txt"
}

func HashHandler(w http.ResponseWriter, r *http.Request) {
	log_v, err := os.ReadFile(filePath())
	if err != nil {
		log.Println(err)
		http.Error(w, "log file not available yet", http.StatusServiceUnavailable)
		return
	}

	pongs := "0"
	if pong_v, err := os.ReadFile(pongFilePath()); err == nil {
		pongs = strings.TrimSpace(string(pong_v))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strings.TrimSpace(string(log_v)) + ".\nPing / Pongs: " + pongs + "\n"))
}

func PeriodicallyPrintHash(hash string) {
	for {
		v := time.Now().UTC().Format(time.RFC3339Nano) + ": " + hash
		fmt.Printf("%s\n", v)

		err := os.WriteFile(filePath(), []byte(v), 0644)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(5 * time.Second)
	}
}

func GenerateRandomHash(n int) error {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return err
	}

	sum := sha256.Sum256(b)
	hash = hex.EncodeToString(sum[:])
	return nil
}

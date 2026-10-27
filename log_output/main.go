package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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

func fetchPongs() (string, error) {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://pingpong-svc:2345/pingpong/pings")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ping_pong returned %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)), nil
}

func HashHandler(w http.ResponseWriter, r *http.Request) {
	log_v, err := os.ReadFile(filePath())
	if err != nil {
		log.Println(err)
		http.Error(w, "log file not available yet", http.StatusServiceUnavailable)
		return
	}

	pongs, err := fetchPongs()
	if err != nil {
		log.Println(err)
		pongs = "0"
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

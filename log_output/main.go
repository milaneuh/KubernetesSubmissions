package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var hash string

func main() {
	port := os.Getenv("PORT")

	err := GenerateRandomHash(512)
	if err != nil {
		log.Fatal(err)
	}

	go PeriodicallyPrintHash(hash)

	http.HandleFunc("/", HashHandler)

	log.Println("Server started in port " + port)
	http.ListenAndServe(":"+port, nil)
}

func HashHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(time.Now().UTC().Format(time.RFC3339Nano) + ": " + hash))
}

func PeriodicallyPrintHash(hash string) {
	for {
		fmt.Printf("%s: %s\n", time.Now().UTC().Format(time.RFC3339Nano), hash)
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

package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

func main() {
	hash, err := GenerateRandomHash(512)
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Printf("%s: %s\n", time.Now().UTC().Format(time.RFC3339Nano), hash)
		time.Sleep(5 * time.Second)
	}
}

func GenerateRandomHash(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

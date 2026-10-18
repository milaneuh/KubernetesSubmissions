package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func main() {
	PrintHash()
	for true {
		time.Sleep(5 * time.Second)
		PrintHash()
	}
}

func PrintHash() {
	hash, err := GenerateRandomHash(512)
	if err != nil {
		panic("Unexpected error when creating a random hash")
	}

	fmt.Printf("%s: %s\n", time.Now().UTC().Format(time.RFC3339Nano), hash)
}
func GenerateRandomHash(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write(b)
	sha := hex.EncodeToString(hasher.Sum(nil))

	return sha, nil
}

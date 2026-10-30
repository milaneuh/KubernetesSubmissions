package main

import (
	"log"
	"os"
	"time"
)

type config struct {
	port                 string
	todoBackendURL       string
	imagePath            string
	imageURL             string
	imageMaxAge          time.Duration
	imageRefreshInterval time.Duration
}

func loadConfig() config {
	return config{
		port:                 mustEnv("PORT"),
		todoBackendURL:       mustEnv("TODO_BACKEND_URL"),
		imagePath:            mustEnv("IMAGE_PATH"),
		imageURL:             mustEnv("IMAGE_URL"),
		imageMaxAge:          mustDuration("IMAGE_MAX_AGE"),
		imageRefreshInterval: mustDuration("IMAGE_REFRESH_INTERVAL"),
	}
}

func mustEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalf("missing required environment variable %s", name)
	}
	return v
}

func mustDuration(name string) time.Duration {
	d, err := time.ParseDuration(mustEnv(name))
	if err != nil {
		log.Fatalf("invalid duration in %s: %v", name, err)
	}
	return d
}

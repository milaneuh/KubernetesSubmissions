package main

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"

	"pingpong/db"
)

//go:embed schema.sql
var schema string

type application struct {
	queries *db.Queries
}

func main() {
	port := mustEnv("PORT")
	databaseURL := mustEnv("DATABASE_URL")

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if _, err := pool.Exec(ctx, schema); err != nil {
		log.Fatal("could not create the schema: ", err)
	}

	queries := db.New(pool)
	if err := queries.InitCounter(ctx); err != nil {
		log.Fatal("could not initialize the counter: ", err)
	}

	app := &application{
		queries: queries,
	}

	http.HandleFunc("/pingpong", app.PingHandler)
	http.HandleFunc("/pingpong/pings", app.PingsHandler)

	log.Println("Server started in port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func mustEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalf("missing required environment variable %s", name)
	}
	return v
}

func (app *application) PingHandler(w http.ResponseWriter, r *http.Request) {
	counter, err := app.queries.IncrementCounter(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong " + strconv.FormatInt(counter, 10)))
}

// PingsHandler returns the current count without incrementing it.
func (app *application) PingsHandler(w http.ResponseWriter, r *http.Request) {
	counter, err := app.queries.GetCounter(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatInt(counter, 10)))
}

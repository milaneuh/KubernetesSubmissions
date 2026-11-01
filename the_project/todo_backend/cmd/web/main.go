package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"todo_backend/internal/db"
	"todo_backend/internal/modules/todo"
)

type Todos interface {
	CreateTodo(ctx context.Context, todo string) error
	ListTodos(ctx context.Context) ([]string, error)
}

type application struct {
	logger *slog.Logger
	todos  Todos
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("missing required environment variable PORT")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if _, err := pool.Exec(ctx, db.Schema); err != nil {
		log.Fatal("could not create the schema: ", err)
	}

	app := &application{
		logger: logger,
		todos:  todo.New(db.New(pool)),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", app.healthz)
	mux.HandleFunc("GET /api/todos", app.listTodos)
	mux.HandleFunc("POST /api/todos", app.newTodo)

	log.Println("Server started in port " + port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("Could not start the server\n", err)
	}

}

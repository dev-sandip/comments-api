package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/dev-sandip/comments-api/internals/env"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	cfg := config{
		address: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost port=5432 user=postgres password=postgres dbname=comments sslmode=disable"),
		},
	}

	//database connection

	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to database")
	defer conn.Close(ctx)
	api := application{config: cfg, db: conn}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	if err := api.run(api.mount()); err != nil {
		slog.Error("Error starting server", "error", err)
	}

}

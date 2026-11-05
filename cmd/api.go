package main

import (
	"log"
	"net/http"
	"time"

	repo "github.com/dev-sandip/comments-api/internals/adapters/postgres/sqlc"
	"github.com/dev-sandip/comments-api/internals/auth"
	"github.com/dev-sandip/comments-api/internals/utils"

	_ "github.com/dev-sandip/comments-api/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config config
	db     *pgx.Conn
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	authService := auth.NewService(repo.New(app.db))
	authHandler := auth.NewHandler(authService)

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Auth routes
	r.Post("/signup", authHandler.SignUpHandler)
	r.Post("/login", authHandler.LoginHandler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{"message": "Hello World!"}
		utils.WriteJSON(w, http.StatusOK, data)
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.address,
		Handler:      h,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Starting server on %s", app.config.address)
	return srv.ListenAndServe()
}

type config struct {
	address string
	db      dbConfig
}

type dbConfig struct {
	dsn string
}

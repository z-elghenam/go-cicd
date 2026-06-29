package main

import (
        "log"
        "net/http"
        "os"
        "time"

        "github.com/bootdotdev/learn-cicd-starter/internal/database"
        "github.com/go-chi/chi/v5"
        "github.com/go-chi/chi/v5/middleware"
        "github.com/joho/godotenv"
        _ "github.com/tursodatabase/libsql-client-go/libsql"
)

type apiConfig struct {
        DB *database.Queries
}

func main() {
        if err := godotenv.Load(); err != nil {
                log.Println("warning: assuming default configuration. .env unreadable:", err)
        }

        port := os.Getenv("PORT")
        if port == "" {
                log.Fatal("PORT environment variable is not set")
        }

        dbURL := os.Getenv("DATABASE_URL")
        if dbURL == "" {
                log.Println("warning: DATABASE_URL not set, using in-memory database")
        }

        db, err := database.NewDB(dbURL)
        if err != nil {
                log.Fatal("Failed to connect to database:", err)
        }
        dbQueries := database.New(db)

        apiCfg := &apiConfig{
                DB: dbQueries,
        }

        r := chi.NewRouter()
        r.Use(middleware.Logger)
        r.Use(middleware.Recoverer)
        r.Use(middleware.Timeout(60 * time.Second))

        fileServer := http.FileServer(http.Dir("./static"))
        r.Handle("/static/*", http.StripPrefix("/static", fileServer))

        r.Get("/healthz", handlerReadiness)

        r.Post("/api/users", apiCfg.handlerUsersCreate)

        r.Group(func(r chi.Router) {
                r.Use(apiCfg.middlewareAuth)
                r.Get("/api/users", apiCfg.handlerUsersGet)
                r.Get("/api/notes", apiCfg.handlerNotesGet)
                r.Post("/api/notes", apiCfg.handlerNotesCreate)
        })

        r.Get("/", func(w http.ResponseWriter, r *http.Request) {
                http.ServeFile(w, r, "./static/index.html")
        })

        server := &http.Server{
                Addr:              ":" + port,
                Handler:           r,
                ReadTimeout:       30 * time.Second,
                WriteTimeout:      30 * time.Second,
                IdleTimeout:       60 * time.Second,
                ReadHeaderTimeout: 10 * time.Second,
        }

        log.Printf("Serving on port: %s\n", port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                log.Fatal("Server failed to start:", err)
        }
}

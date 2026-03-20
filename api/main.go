package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kraken/api/app"
	"kraken/api/repository"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := mustOpenDB("./example.db")
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	sensorsRepo := repository.NewSensorsRepository(db)
	app := app.NewApp(userRepo, sensorsRepo)

	server := &http.Server{
		Addr:         ":12345",
		Handler:      app.Router(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Println("Server running on :12345")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	gracefulShutdown(server)
}

func mustOpenDB(dsn string) *sql.DB {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// SQLite pool configuration
	db.SetMaxOpenConns(1) // important for SQLite
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	// Check connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

func gracefulShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped")
}

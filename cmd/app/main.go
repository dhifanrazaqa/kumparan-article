package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dhifanrazaqa/kumparan-article/internal/handlers"
	"github.com/dhifanrazaqa/kumparan-article/internal/repositories"
	"github.com/dhifanrazaqa/kumparan-article/internal/router"
	"github.com/dhifanrazaqa/kumparan-article/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	loadEnv()

	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v\n", err)
	}
	defer dbPool.Close()
	log.Println("Connected to the database successfully.")

	userRepo := repositories.NewPgxUserRepo(dbPool)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	routerDeps := router.Deps{
		UserHandler: userHandler,
	}

	mainRouter := router.SetupRouter(routerDeps)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mainRouter,
	}

	go func() {
		log.Printf("Server started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Received shutdown signal, shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server shut down successfully.")
}

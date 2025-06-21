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
	"github.com/go-redis/redis"
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
	redisURL := os.Getenv("REDIS_URL")
	port := os.Getenv("APP_PORT")
	jwtSecret := os.Getenv("JWT_SECRET")
	refreshTokenSecret := os.Getenv("REFRESH_TOKEN_SECRET")

	if port == "" {
		port = "8080"
	}
	if dbURL == "" || redisURL == "" {
		log.Fatal("Error: DATABASE_URL dan REDIS_URL harus diatur")
	}

	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v\n", err)
	}
	defer dbPool.Close()

	redisClient := redis.NewClient(&redis.Options{Addr: redisURL})
	if _, err := redisClient.Ping().Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Database and Redis.")

	userRepo := repositories.NewPgxUserRepo(dbPool)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	authService := services.NewAuthService(userRepo, jwtSecret, refreshTokenSecret, redisClient)
	authHandler := handlers.NewAuthHandler(authService)

	routerDeps := router.Deps{
		UserHandler: userHandler,
		AuthHandler: authHandler,
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

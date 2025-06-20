package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		log.Fatalf("Tidak dapat terhubung ke database: %v\n", err)
	}
	defer dbPool.Close()

	mainRouter := http.NewServeMux()
	mainRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Selamat datang di aplikasi Go!"))
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mainRouter,
	}

	go func() {
		log.Printf("Server mulai di port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Gagal memulai server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Menerima sinyal shutdown, mematikan server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown gagal: %v", err)
	}
	log.Println("Server berhasil dimatikan.")
}

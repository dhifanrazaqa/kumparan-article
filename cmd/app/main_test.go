package main

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/dhifanrazaqa/kumparan-article/internal/handlers"
	"github.com/dhifanrazaqa/kumparan-article/internal/repositories"
	"github.com/dhifanrazaqa/kumparan-article/internal/router"
	"github.com/dhifanrazaqa/kumparan-article/internal/services"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var testRouter *mux.Router
var testDbPool *pgxpool.Pool

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Peringatan: Gagal memuat file .env")
	}

	testDbURL := os.Getenv("TEST_DATABASE_URL")
	if testDbURL == "" {
		testDbURL = "postgres://user:password@localhost:5433/articledb_test?sslmode=disable"
	}
	var err error
	testDbPool, err = pgxpool.New(context.Background(), testDbURL)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database tes: %v", err)
	}
	defer testDbPool.Close()

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{Addr: redisURL})
	if _, err := redisClient.Ping().Result(); err != nil {
		log.Fatalf("Gagal terhubung ke Redis tes: %v", err)
	}
	log.Println("Database dan Redis untuk tes berhasil terhubung.")

	clearDatabase(testDbPool)

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	refreshTokenSecret := os.Getenv("REFRESH_TOKEN_SECRET")

	userRepo := repositories.NewPgxUserRepo(testDbPool)
	articleRepo := repositories.NewPgxArticleRepo(testDbPool)

	authService := services.NewAuthService(userRepo, jwtSecret, refreshTokenSecret, redisClient)
	userService := services.NewUserService(userRepo)
	articleService := services.NewArticleService(articleRepo, redisClient)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	articleHandler := handlers.NewArticleHandler(articleService)

	routerDeps := router.Deps{
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		ArticleHandler: articleHandler,
		JWTSecret:      jwtSecret,
	}
	testRouter = router.SetupRouter(routerDeps)

	exitCode := m.Run()

	testDbPool.Close()
	os.Exit(exitCode)
}

func clearDatabase(pool *pgxpool.Pool) {
	_, err := pool.Exec(context.Background(), "TRUNCATE TABLE articles, users RESTART IDENTITY CASCADE")
	if err != nil {
		log.Fatalf("Gagal membersihkan database: %v", err)
	}
}

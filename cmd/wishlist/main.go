package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wishlist/internal/application"
	"wishlist/internal/config"
	"wishlist/internal/infrastructure/in"
	"wishlist/internal/infrastructure/in/handlers"
	"wishlist/internal/infrastructure/in/middleware"
	"wishlist/internal/infrastructure/out"
	"wishlist/internal/infrastructure/out/repos"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m, err := migrate.New("file://migrations", cfg.GetDSN())
	if err != nil {
		log.Fatalf("create migrations: %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("apply migrations: %v", err)
	}

	dbpool, err := pgxpool.New(ctx, cfg.GetDSN())
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer dbpool.Close()

	userRepo := repos.NewUserRepo(dbpool)
	wishlistRepo := repos.NewWishlistRepo(dbpool)
	itemRepo := repos.NewItemRepo(dbpool)

	hasher := out.NewBcryptHasher()
	tokenManager := out.NewJWTManager(cfg.JWTSecret, 24*time.Hour)
	publicTokenGenerator := out.NewPublicTokenGenerator()

	authService := application.NewAuthService(userRepo, hasher, tokenManager)
	wishlistService := application.NewWishlistService(wishlistRepo, publicTokenGenerator)
	itemService := application.NewItemService(itemRepo, wishlistRepo)

	authHandler := handlers.NewAuthHandler(authService)
	wishlistHandler := handlers.NewWishlistHandler(wishlistService, itemService)
	itemHandler := handlers.NewItemHandler(itemService, wishlistService)

	authMiddleware := middleware.NewAuthMiddleware(tokenManager)

	router := in.NewRouter(authHandler, wishlistHandler, itemHandler, authMiddleware)

	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	go func() {
		log.Printf("server started on %s", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down server...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}

	log.Println("server stopped gracefully")
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vadimpk/ppc-project/controller"
	"github.com/vadimpk/ppc-project/controller/middleware"
	"github.com/vadimpk/ppc-project/pkg/auth"
	"github.com/vadimpk/ppc-project/repository"
	"github.com/vadimpk/ppc-project/services"
)

const (
	defaultPort            = "8080"
	defaultShutdownTimeout = 10 * time.Second
)

func main() {
	// Initialize database connection
	db, err := repository.NewDB(repository.Options{
		Host:           "localhost",
		Port:           5432,
		User:           "postgres",
		Pass:           "postgres",
		DBName:         "ppc",
		MinConnections: 1,
		MaxConnections: 2,
		Timezone:       "UTC",
	})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	repositories := repository.NewRepositories(db)

	// Initialize JWT token manager
	tokenManager, err := auth.NewTokenManager("your-signing-key")
	if err != nil {
		log.Fatalf("Failed to initialize token manager: %v", err)
	}

	// Initialize services
	srvcs := services.NewServices(repositories)

	// Initialize handlers and middleware
	handlers := controller.NewHandlers(srvcs, tokenManager)
	authMiddleware := middleware.NewAuthMiddleware(tokenManager)

	// Initialize router
	router := controller.NewRouter(handlers, authMiddleware.Authenticate)

	// Configure server
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
		// Good practice settings
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	// Kill (no param) default send syscall.SIGTERM
	// Kill -2 is syscall.SIGINT
	// Kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

// Debug helper - remove in production
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

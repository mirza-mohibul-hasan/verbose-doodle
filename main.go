package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"qrcode-generator/internal/handler"
	"qrcode-generator/internal/middleware"
	"qrcode-generator/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Configure structured JSON logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Initialize services
	qrService := service.NewQRService()
	qrHandler := handler.NewQRHandler(qrService)

	// Initialize rate limiter: 10 requests/second, burst of 20
	rateLimiter := middleware.NewRateLimiter(10, 20)

	// Setup router
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Serve static files (frontend)
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	// API routes
	api := router.Group("/api")
	api.Use(rateLimiter.Middleware())
	{
		api.POST("/generate", qrHandler.Generate)
		api.GET("/health", qrHandler.Health)
	}

	// Configure HTTP server with timeouts
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		slog.Info("Server starting", "port", port, "url", "http://localhost:"+port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited gracefully")
}

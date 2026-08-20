package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/article-management/backend/internal/article"
	"github.com/example/article-management/backend/internal/config"
	"github.com/example/article-management/backend/internal/database"
	appmiddleware "github.com/example/article-management/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("load configuration", "error", err)
		os.Exit(1)
	}
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.Open(cfg)
	if err != nil {
		logger.Error("connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	repository := article.NewSQLRepository(db)
	service := article.NewService(repository)
	handler := article.NewHandler(service, logger)

	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger(), appmiddleware.CORS(cfg.CORSAllowedOrigins))
	router.GET("/health", func(c *gin.Context) {
		if err := db.PingContext(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	article.RegisterRoutes(router, handler)

	server := &http.Server{Addr: cfg.Address(), Handler: router, ReadHeaderTimeout: cfg.ShutdownTimeout}
	go func() {
		logger.Info("API server started", "address", cfg.Address())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("serve API", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("shutdown API", "error", err)
	}
}

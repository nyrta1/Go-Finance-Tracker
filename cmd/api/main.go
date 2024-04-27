package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-finance-tracker/internal/config"
	"go-finance-tracker/internal/db/psql"
	"go-finance-tracker/internal/repository"
	"go-finance-tracker/internal/rest/handler"
	"go-finance-tracker/internal/rest/routers"
	"go-finance-tracker/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var appConfig config.App

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	logger.InitLogger()

	appConfig = config.App{
		PORT: os.Getenv("APP_PORT"),
		DB:   initializeDB(),
	}

	dbInstance, err := psql.GetDbInstance(appConfig.DB)
	if err != nil {
		logger.GetLogger().Fatal("Error initializing DB:", err)
	}

	userRepo := repository.NewUserRepository(dbInstance)
	roleRepo := repository.NewRoleRepository(dbInstance)

	authHandlers := handler.NewAuthHandler(userRepo, roleRepo)

	r := gin.Default()

	router := routers.NewRouters(*authHandlers)
	router.SetupRoutes(r)
	r.Use(rateLimitMiddleware())

	server := &http.Server{
		Addr:    ":" + appConfig.PORT,
		Handler: r,
	}

	gracefulShutdown(server)
}

func initializeDB() config.PostgresDB {
	dbConfig := config.PostgresDB{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Sslmode:  os.Getenv("POSTGRES_SSLMODE"),
		Name:     os.Getenv("POSTGRES_NAME"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}

	return dbConfig
}

func gracefulShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		logger.GetLogger().Info("Server is shutting down...")

		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.GetLogger().Fatal("Server shutdown error:", err)
		}

		logger.GetLogger().Info("Server has gracefully stopped")
		os.Exit(0)
	}()

	logger.GetLogger().Info("Server is running on :" + appConfig.PORT)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.GetLogger().Fatal("Error starting server:", err)
	}
}

func rateLimitMiddleware() gin.HandlerFunc {
	limiter := time.Tick(time.Second)

	return func(c *gin.Context) {
		select {
		case <-limiter:
			c.Next()
		default:
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
		}
	}
}

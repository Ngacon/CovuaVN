package main

import (  //bố m tự viết nhé
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"covuavn/internal/article"
	"covuavn/internal/auth"
	"covuavn/internal/chat"
	"covuavn/internal/config"
	"covuavn/internal/db"
	"covuavn/internal/notification"
	"covuavn/internal/router"
	"covuavn/internal/user"
)

func main() {
	cfg := config.Load()

	pool, err := db.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer pool.Close()

	userRepo := user.NewRepository(pool)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	authRepo := auth.NewRepository(pool)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	chatRepo := chat.NewRepository(pool)
	chatService := chat.NewService(chatRepo)
	chatHandler := chat.NewHandler(chatService)

	notificationRepo := notification.NewRepository(pool)
	notificationService := notification.NewService(notificationRepo)
	notificationHandler := notification.NewHandler(notificationService)

	articleRepo := article.NewRepository(pool)
	articleService := article.NewService(articleRepo)
	articleHandler := article.NewHandler(articleService)

	r := router.New(userHandler, authHandler, chatHandler, notificationHandler, articleHandler)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("api listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

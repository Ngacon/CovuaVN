package main

import (
	"database/sql"
	"os"

	"github.com/Ngacon/covuavn/internal/delivery"
	"github.com/Ngacon/covuavn/internal/repository"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, _ := sql.Open("postgres", os.Getenv("DB_URL"))

	userRepo := &repository.UserRepo{DB: db}
	authHand := &delivery.AuthHandler{Repo: userRepo}

	r := gin.Default()
	r.POST("/register", authHand.Register)
	r.Run()
}

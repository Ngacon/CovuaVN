package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/Ngacon/covuavn/internal/domain"
	"github.com/Ngacon/covuavn/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Repo *repository.UserRepo
}

func (h *AuthHandler) Register(c *gin.Context) {
	var u domain.User
	c.BindJSON(&u)
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	u.Password = string(hash)
	h.Repo.Create(u)
	c.JSON(200, gin.H{"status": "ok"})
}
package v1

import (
	"gin-user-management/internal/repository"

	"github.com/gin-gonic/gin"
)

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (as *authService) Login(ctx *gin.Context, email, password string) error {
	return nil
}

func (as *authService) Logout(ctx *gin.Context) error {
	return nil
}

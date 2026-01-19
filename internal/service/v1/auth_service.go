package v1

import (
	"gin-user-management/internal/repository"
	"gin-user-management/internal/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	context := ctx.Request.Context()

	email = util.NormalizeString(email)
	user, err := as.repo.GetByEmail(context, email)
	if err != nil {
		return util.NewError("Invalid email or password.", util.ErrCodeUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return util.NewError("Invalid email or password.", util.ErrCodeUnauthorized)
	}

	return nil
}

func (as *authService) Logout(ctx *gin.Context) error {
	return nil
}

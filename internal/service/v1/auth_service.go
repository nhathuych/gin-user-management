package v1

import (
	"gin-user-management/internal/repository"
	"gin-user-management/internal/util"
	"gin-user-management/pkg/auth"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo           repository.UserRepository
	tokenGenerator auth.TokenGenerator
}

func NewAuthService(repo repository.UserRepository, tokenGenerator auth.TokenGenerator) AuthService {
	return &authService{
		repo:           repo,
		tokenGenerator: tokenGenerator,
	}
}

func (as *authService) Login(ctx *gin.Context, email, password string) (string, string, int, error) {
	context := ctx.Request.Context()

	email = util.NormalizeString(email)
	user, err := as.repo.GetByEmail(context, email)
	if err != nil {
		return "", "", 0, util.NewError("Invalid credentials.", util.ErrCodeUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", 0, util.NewError("Invalid credentials.", util.ErrCodeUnauthorized)
	}

	accessToken, err := as.tokenGenerator.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, util.WrapError(err, "Failed to generate access token.", util.ErrCodeInternal)
	}

	refreshToken, err := as.tokenGenerator.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, util.WrapError(err, "Failed to generate refresh token.", util.ErrCodeInternal)
	}

	if err := as.tokenGenerator.StoreRefreshToken(context, refreshToken); err != nil {
		return "", "", 0, util.WrapError(err, "Failed to store refresh token.", util.ErrCodeInternal)
	}

	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}

func (as *authService) Logout(ctx *gin.Context) error {
	return nil
}

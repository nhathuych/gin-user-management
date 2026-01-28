package v1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gin-user-management/internal/db/sqlc"
	"gin-user-management/internal/repository"
	"gin-user-management/internal/util"
	"gin-user-management/pkg/cache"
	"gin-user-management/pkg/logger"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo       repository.UserRepository
	redisCache cache.RedisCacheService
}

func NewUserService(repo repository.UserRepository, redisCache cache.RedisCacheService) UserService {
	return &userService{
		repo:       repo,
		redisCache: redisCache,
	}
}

func (us *userService) GetAll(ctx *gin.Context, search, orderBy, sort string, page, limit int32, deleted bool) ([]sqlc.User, int32, error) {
	context := ctx.Request.Context()

	search = strings.TrimSpace(search)
	if sort == "" {
		sort = "desc"
	}
	if orderBy == "" {
		orderBy = "id"
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var cacheData cache.PagedResult[sqlc.User]
	cacheKey := us.userListCacheKey(search, orderBy, sort, page, limit, deleted)
	if err := us.redisCache.Get(context, cacheKey, &cacheData); err == nil {
		return cacheData.Items, cacheData.Total, nil
	}

	offset := (page - 1) * limit

	users, err := us.repo.GetAllV2(context, search, orderBy, sort, limit, offset, deleted)
	if err != nil {
		return []sqlc.User{}, 0, util.WrapError(err, "Failed to retrieve users.", util.ErrCodeInternal)
	}

	total, err := us.repo.CountUsers(context, search, deleted)
	if err != nil {
		return []sqlc.User{}, 0, util.WrapError(err, "Failed to count users.", util.ErrCodeInternal)
	}

	cacheData = cache.PagedResult[sqlc.User]{
		Items: users,
		Total: int32(total),
	}
	us.redisCache.Set(context, cacheKey, cacheData, 5*time.Minute)

	return users, int32(total), nil
}

func (us *userService) Create(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, util.WrapError(err, "Failed to hash password.", util.ErrCodeInternal)
	}

	input.Password = string(hashedPassword)
	input.Email = util.NormalizeString(input.Email)

	user, err := us.repo.Create(context, input)
	if err != nil {
		return sqlc.User{}, util.WrapError(err, "Failed to create a new user.", util.ErrCodeInternal)
	}

	us.clearUserListCache(context)

	return user, nil
}

func (us *userService) GetByUUID(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()

	user, err := us.repo.GetByUUID(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, util.NewError("User not found.", util.ErrCodeNotFound)
		}
		return sqlc.User{}, util.WrapError(err, "Failed to retrieve user.", util.ErrCodeInternal)
	}

	return user, nil
}

func (us *userService) Update(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	if input.Password != nil && *input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return sqlc.User{}, util.WrapError(err, "Failed to hash password.", util.ErrCodeInternal)
		}

		strPassword := string(hashedPassword)
		input.Password = &strPassword
	}

	user, err := us.repo.Update(context, input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, util.NewError("User not found.", util.ErrCodeNotFound)
		}
		return sqlc.User{}, util.WrapError(err, "Failed to update user.", util.ErrCodeInternal)
	}

	us.clearUserListCache(context)

	return user, nil
}

func (us *userService) SoftDeleteUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()

	user, err := us.repo.SoftDeleteUser(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, util.NewError("User not found.", util.ErrCodeNotFound)
		}
		return sqlc.User{}, util.WrapError(err, "Failed to delete user.", util.ErrCodeInternal)
	}

	us.clearUserListCache(context)

	return user, nil
}

func (us *userService) RestoreUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()

	user, err := us.repo.RestoreUser(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, util.NewError("User not found or not deleted yet.", util.ErrCodeNotFound)
		}
		return sqlc.User{}, util.WrapError(err, "Failed to restore user.", util.ErrCodeInternal)
	}

	us.clearUserListCache(context)

	return user, nil
}

func (us *userService) HardDeleteUser(ctx *gin.Context, uuid uuid.UUID) error {
	context := ctx.Request.Context()

	_, err := us.repo.HardDeleteUser(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return util.NewError("User not found or not eligible for hard delete.", util.ErrCodeNotFound)
		}
		return util.WrapError(err, "Failed to delete user.", util.ErrCodeInternal)
	}

	us.clearUserListCache(context)

	return nil
}

func (us *userService) userListCacheKey(search, orderBy, sort string, page, limit int32, deleted bool) string {
	return fmt.Sprintf("users:list:v1:%s:%s:%s:%d:%d:%t", search, orderBy, sort, page, limit, deleted)
}

func (us *userService) clearUserListCache(ctx context.Context) {
	if err := us.redisCache.Clear(ctx, "users:list:v1:*"); err != nil {
		logger.AppLogger.Warn().Err(err).Msg("Failed to clear users list cache")
	}
}

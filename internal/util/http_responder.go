package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseSuccess[T any](ctx *gin.Context, status int, message string, data T, pagination *Pagination) {
	ctx.JSON(status, APIResponse[T]{
		Status:     "success",
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func ResponseError(ctx *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		ctx.JSON(httpStatusFromCode(appErr.Code), gin.H{
			"code":  appErr.Code,
			"error": appErr.Message,
		})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"code":  ErrCodeInternal,
		"error": err.Error(),
	})
}

func httpStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func ResponseStatusCode(ctx *gin.Context, status int) {
	ctx.Status(status)
}

func ResponseValidator(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusBadRequest, data)
}

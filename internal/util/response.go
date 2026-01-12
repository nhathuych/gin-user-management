package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

const (
	ErrCodeBadRequest ErrorCode = "BAD_REQUEST"
	ErrCodeNotFound   ErrorCode = "NOT_FOUND"
	ErrCodeConflict   ErrorCode = "CONFLICT"
	ErrCodeInternal   ErrorCode = "INTERNAL_SERVER_ERROR"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (ae *AppError) Error() string {
	return ""
}

func NewError(message string, code ErrorCode) error {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func WrapError(err error, message string, code ErrorCode) error {
	return &AppError{
		Code:    code,
		Err:     err,
		Message: message,
	}
}

func ResponseError(ctx *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		status := httpStatusFromCode(appErr.Code)
		response := gin.H{
			"code":  appErr.Code,
			"error": appErr.Message,
		}

		if appErr.Err != nil {
			response["detail"] = appErr.Err.Error()
		}

		ctx.JSON(status, response)
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"code":  ErrCodeInternal,
		"error": err.Error(),
	})
}

func ResponseSuccess(ctx *gin.Context, status int, message string, data ...any) {
	res := APIResponse{
		Status:  "success",
		Message: message,
	}
	if len(data) > 0 && data[0] != nil {
		res.Data = data[0]
	}

	ctx.JSON(status, res)
}

func ResponseStatusCode(ctx *gin.Context, status int) {
	ctx.Status(status)
}

func ResponseValidator(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusBadRequest, data)
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

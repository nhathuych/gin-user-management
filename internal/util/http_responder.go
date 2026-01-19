package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Detail  string    `json:"detail,omitempty"`
}

type ErrorResponse struct {
	Status string      `json:"status"`
	Error  ErrorDetail `json:"error"`
}

type APIResponse[T any] struct {
	Status     string      `json:"status"`
	Message    string      `json:"message,omitempty"`
	Data       T           `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func ResponseSuccess[T any](ctx *gin.Context, status int, message string, data T, pagination *Pagination) {
	ctx.JSON(status, APIResponse[T]{
		Status:     "success",
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func ResponseMessage(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, APIResponse[any]{
		Status:  "success",
		Message: message,
	})
}

func ResponseError(ctx *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		status := httpStatusFromCode(appErr.Code)

		res := ErrorResponse{
			Status: "error",
			Error: ErrorDetail{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		}

		if appErr.Err != nil {
			res.Error.Detail = appErr.Err.Error()
		}

		ctx.JSON(status, res)
		return
	}

	ctx.JSON(http.StatusInternalServerError, ErrorResponse{
		Status: "error",
		Error: ErrorDetail{
			Code:    ErrCodeInternal,
			Message: "Internal server error",
			Detail:  err.Error(),
		},
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
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
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

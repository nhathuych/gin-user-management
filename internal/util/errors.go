package util

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

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
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
		Message: message,
		Err:     err,
	}
}

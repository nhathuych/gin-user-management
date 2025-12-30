package validation

import (
	"fmt"
	"gin-user-management/internal/util"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to get validator engine")
	}

	RegisterCustomValidation(v)
	return nil
}

func HandleValidationErrors(err error) gin.H {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)

		for _, e := range validationError {
			root := strings.Split(e.Namespace(), ".")[0]
			rawPath := strings.TrimPrefix(e.Namespace(), root+".")
			parts := strings.Split(rawPath, ".")

			for i, part := range parts {
				if strings.Contains(part, "[") {
					idx := strings.Index(part, "[")
					base := util.CamelToSnake(part[:idx])
					index := part[idx:]
					parts[i] = base + index
				} else {
					parts[i] = util.CamelToSnake(part)
				}
			}

			fieldPath := strings.Join(parts, ".")

			switch e.Tag() {
			case "gt":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than %s", fieldPath, e.Param())
			case "lt":
				errors[fieldPath] = fmt.Sprintf("%s must be less than %s", fieldPath, e.Param())
			case "gte":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than or equal to %s", fieldPath, e.Param())
			case "lte":
				errors[fieldPath] = fmt.Sprintf("%s must be less than or equal to %s", fieldPath, e.Param())
			case "uuid":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid UUID", fieldPath)
			case "slug":
				errors[fieldPath] = fmt.Sprintf("%s can only contain lowercase letters, numbers, hyphens or dots", fieldPath)
			case "min":
				errors[fieldPath] = fmt.Sprintf("%s must be at least %s characters", fieldPath, e.Param())
			case "max":
				errors[fieldPath] = fmt.Sprintf("%s must be at most %s characters", fieldPath, e.Param())
			case "min_int":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than %s", fieldPath, e.Param())
			case "max_int":
				errors[fieldPath] = fmt.Sprintf("%s must be less than %s", fieldPath, e.Param())
			case "oneof":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ", ")
				errors[fieldPath] = fmt.Sprintf("%s must be one of the following values: %s", fieldPath, allowedValues)
			case "required":
				errors[fieldPath] = fmt.Sprintf("%s is required", fieldPath)
			case "search":
				errors[fieldPath] = fmt.Sprintf("%s can only contain letters (upper or lower case), numbers and spaces", fieldPath)
			case "email":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid email address", fieldPath)
			case "datetime":
				errors[fieldPath] = fmt.Sprintf("%s must be in the format YYYY-MM-DD", fieldPath)
			case "file_ext":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ", ")
				errors[fieldPath] = fmt.Sprintf("%s only allows files with the following extensions: %s", fieldPath, allowedValues)
			}
		}

		return gin.H{"error": errors}
	}

	return gin.H{"error": err.Error()}
}

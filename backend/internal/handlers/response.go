package handlers

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field string `json:"field"`
	Rule  string `json:"rule"`
}

func JSONError(c *gin.Context, status int, code string, message string) {
	payload := gin.H{
		"error": message,
		"code":  code,
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		payload["request_id"] = requestID
	}

	c.JSON(status, payload)
}

func JSONValidationError(c *gin.Context, message string, err error) {
	payload := gin.H{
		"error": message,
		"code":  "VALIDATION_ERROR",
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		payload["request_id"] = requestID
	}

	if fields := validationFields(err); len(fields) > 0 {
		payload["fields"] = fields
	}

	c.JSON(400, payload)
}

func RecordError(c *gin.Context, err error) {
	if err != nil {
		_ = c.Error(err)
	}
}

func validationFields(err error) []FieldError {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return nil
	}

	fields := make([]FieldError, 0, len(validationErrors))
	for _, fieldErr := range validationErrors {
		fields = append(fields, FieldError{
			Field: jsonFieldName(fieldErr),
			Rule:  fieldErr.Tag(),
		})
	}

	return fields
}

func jsonFieldName(fieldErr validator.FieldError) string {
	field := fieldErr.Field()
	if field == "" {
		return ""
	}

	return toSnakeCase(field)
}

func toSnakeCase(value string) string {
	var builder strings.Builder
	for index, char := range value {
		if index > 0 && char >= 'A' && char <= 'Z' {
			builder.WriteByte('_')
		}

		builder.WriteRune(char)
	}

	return strings.ToLower(builder.String())
}

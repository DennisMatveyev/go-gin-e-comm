package common

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateRequest[T any](c *gin.Context, model *T) error {
	if err := c.ShouldBindJSON(model); err != nil {
		return ErrParseJSON
	}
	if err := validate.Struct(model); err != nil {
		return err
	}
	return nil
}

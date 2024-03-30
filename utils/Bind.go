package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{validator: validator.New()}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func BindBody(
	c *fiber.Ctx, input interface{},
) error {
	if err := c.BodyParser(input); err != nil {
		return err
	}
	if err := NewValidator().Validate(input); err != nil {
		return err
	}

	return nil
}

func BindQueries(
	c *fiber.Ctx, input interface{}, queriesList []string,
) error {
	if err := c.QueryParser(input); err != nil {
		return err
	}
	if err := NewValidator().Validate(input); err != nil {
		return err
	}

	return nil
}

package utils

import (
	"github.com/JohnKucharsky/DevGroup/domain"
	"github.com/gofiber/fiber/v2"
)

func ErrorRes(errorString string) fiber.Map {
	return fiber.Map{
		"success": false,
		"message": errorString,
	}
}

func SuccessRes(data interface{}) fiber.Map {
	return fiber.Map{
		"success": true,
		"data":    data,
	}
}

func SuccessPaginatedRes(data interface{}, pagination *domain.Pagination) fiber.Map {
	return fiber.Map{
		"success":    true,
		"data":       data,
		"pagination": pagination,
	}
}

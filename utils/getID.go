package utils

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetID(c *fiber.Ctx) (int, error) {
	var id = c.Params("id")
	if id == "" {
		return 0, errors.New("no id in the url")
	}
	parsedID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, err
	}

	return int(parsedID), nil
}

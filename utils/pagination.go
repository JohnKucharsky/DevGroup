package utils

import (
	"github.com/JohnKucharsky/DevGroup/domain"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetPaginationParams(c *fiber.Ctx) (*domain.ParsedPaginationParams, error) {
	var pp domain.PaginationParams
	if err := BindQueries(c, &pp, []string{"limit", "offset"}); err != nil {
		return nil, err
	}

	limit, _ := strconv.Atoi(pp.Limit)
	offset, _ := strconv.Atoi(pp.Offset)

	if limit > 0 {
		if offset > 0 {
			return &domain.ParsedPaginationParams{
				Limit:  limit,
				Offset: &offset,
			}, nil
		}
		return &domain.ParsedPaginationParams{
			Limit:  limit,
			Offset: nil,
		}, nil
	}

	return nil, nil
}

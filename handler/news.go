package handler

import (
	"github.com/JohnKucharsky/DevGroup/domain"
	"github.com/JohnKucharsky/DevGroup/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (h *Handler) CreateNews(c *fiber.Ctx) error {
	var req domain.NewsInput
	if err := utils.BindBody(c, &req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	res, err := h.newsStore.Create(req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(utils.SuccessRes(res))
}

func (h *Handler) GetManyNews(c *fiber.Ctx) error {
	pp, err := utils.GetPaginationParams(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	res, pagination, err := h.newsStore.GetManyPaginated(pp)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(utils.SuccessPaginatedRes(res, pagination))
}

func (h *Handler) UpdateNews(c *fiber.Ctx) error {
	id, err := utils.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	var req domain.NewsInputUpdate
	if err := utils.BindBody(c, &req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	res, err := h.newsStore.Update(req, id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(utils.SuccessRes(res))
}

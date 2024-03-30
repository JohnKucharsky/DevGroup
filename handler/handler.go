package handler

import (
	"github.com/JohnKucharsky/DevGroup/domain"
)

type Handler struct {
	userStore domain.AuthStore
	newsStore domain.NewsStore
}

func NewHandler(
	us domain.AuthStore,
	ns domain.NewsStore,

) *Handler {
	return &Handler{
		userStore: us,
		newsStore: ns,
	}
}

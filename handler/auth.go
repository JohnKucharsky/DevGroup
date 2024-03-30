package handler

import (
	"github.com/JohnKucharsky/DevGroup/domain"
	"github.com/JohnKucharsky/DevGroup/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) SignUp(c *fiber.Ctx) error {
	var req domain.SignUpInput
	if err := utils.BindBody(c, &req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	err := req.HashPassword()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	res, err := h.userStore.Create(req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(utils.SuccessRes(res))
}

func (h *Handler) SignIn(c *fiber.Ctx) error {
	var req domain.SignInInput
	if err := utils.BindBody(c, &req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	res, err := h.userStore.GetOne(strings.ToLower(req.Email), "")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	ok, err := req.CheckPassword(res.Password)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(utils.ErrorRes("passwords don't match"))
	}

	accessToken, err := h.userStore.SetAccessToken(res.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	refreshToken, err := h.userStore.SetRefreshToken(res.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	var accessTokenMaxAgeString = os.Getenv("ACCESS_TOKEN_MAXAGE")
	var refreshTokenMaxAgeString = os.Getenv("REFRESH_TOKEN_MAXAGE")
	var accessTokenMaxAge, _ = strconv.Atoi(accessTokenMaxAgeString)
	var refreshTokenMaxAge, _ = strconv.Atoi(refreshTokenMaxAgeString)

	c.Cookie(
		&fiber.Cookie{
			Name:     "access_token",
			Value:    *accessToken,
			Path:     "/",
			MaxAge:   accessTokenMaxAge * 60,
			Secure:   false,
			HTTPOnly: true,
		},
	)

	c.Cookie(
		&fiber.Cookie{
			Name:     "refresh_token",
			Value:    *refreshToken,
			Path:     "/",
			MaxAge:   refreshTokenMaxAge * 60,
			Secure:   false,
			HTTPOnly: true,
		},
	)

	c.Cookie(
		&fiber.Cookie{
			Name:     "logged_in",
			Value:    "true",
			Path:     "/",
			MaxAge:   accessTokenMaxAge * 60,
			Secure:   false,
			HTTPOnly: false,
		},
	)

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success":      true,
			"access_token": accessToken,
			"user":         res,
		},
	)
}

func (h *Handler) RefreshAccessToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(http.StatusBadRequest).JSON(
			utils.ErrorRes("no refresh token in cookies"),
		)
	}

	userID, err := h.userStore.GetByRefreshTokenRedis(refreshToken)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(utils.ErrorRes(err.Error()))
	}

	res, err := h.userStore.GetOne("", userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	accessToken, err := h.userStore.SetAccessToken(res.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	var accessTokenMaxAgeString = os.Getenv("ACCESS_TOKEN_MAXAGE")
	var accessTokenMaxAge, _ = strconv.Atoi(accessTokenMaxAgeString)

	c.Cookie(
		&fiber.Cookie{
			Name:     "access_token",
			Value:    *accessToken,
			Path:     "/",
			MaxAge:   accessTokenMaxAge * 60,
			Secure:   false,
			HTTPOnly: true,
		},
	)

	c.Cookie(
		&fiber.Cookie{
			Name:     "logged_in",
			Value:    "true",
			Path:     "/",
			MaxAge:   accessTokenMaxAge * 60,
			Secure:   false,
			HTTPOnly: false,
		},
	)

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"success":      true,
			"access_token": accessToken,
		},
	)
}

func (h *Handler) DeserializeUser(c *fiber.Ctx) error {
	var accessToken string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		accessToken = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("access_token") != "" {
		accessToken = c.Cookies("access_token")
	}

	if accessToken == "" {
		return c.Status(http.StatusUnauthorized).JSON(
			utils.ErrorRes("No access token"),
		)
	}

	userID, tokenUUID, err := h.userStore.GetByAccessTokenRedis(accessToken)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(utils.ErrorRes(err.Error()))
	}

	res, err := h.userStore.GetOne("", userID)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(utils.ErrorRes(err.Error()))
	}

	c.Locals("user", res)
	c.Locals("access_token_uuid", tokenUUID)

	return c.Next()
}

func (h *Handler) GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)

	return c.Status(http.StatusOK).JSON(utils.SuccessRes(user))
}

func (h *Handler) LogoutUser(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(http.StatusUnauthorized).JSON(utils.ErrorRes("No refresh token in the cookies"))
	}
	accessToken, ok := c.Locals("access_token_uuid").(string)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(utils.ErrorRes("Access token is not a string"))
	}
	if accessToken == "" {
		return c.Status(http.StatusUnauthorized).JSON(utils.ErrorRes("Access token an empty string"))
	}

	err := h.userStore.DeleteTokensRedis(refreshToken, accessToken)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	now := time.Now()

	c.Cookie(
		&fiber.Cookie{
			Name:    "access_token",
			Value:   "",
			Expires: now,
		},
	)
	c.Cookie(
		&fiber.Cookie{
			Name:    "refresh_token",
			Value:   "",
			Expires: now,
		},
	)
	c.Cookie(
		&fiber.Cookie{
			Name:    "logged_in",
			Value:   "",
			Expires: now,
		},
	)

	return c.SendStatus(http.StatusOK)
}

package router

import (
	"github.com/JohnKucharsky/DevGroup/handler"
	"github.com/JohnKucharsky/DevGroup/store"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Register(r *fiber.App, db *pgxpool.Pool, redis *redis.Client) {
	us := store.NewUserStore(db, redis)
	ns := store.NewNewsStore(db)

	h := handler.NewHandler(
		us,
		ns,
	)

	v1 := r.Group("/api")

	// auth
	auth := v1.Group("/auth")
	auth.Post("/sign-up", h.SignUp)
	auth.Post("/login", h.SignIn)
	auth.Get("/logout", h.DeserializeUser, h.LogoutUser)
	auth.Get("/refresh", h.RefreshAccessToken)
	auth.Get("/me", h.DeserializeUser, h.GetMe)
	// end auth

	// news
	v1.Post("/add", h.DeserializeUser, h.CreateNews)
	v1.Post("/edit/:id", h.DeserializeUser, h.UpdateNews)
	v1.Get("/list", h.DeserializeUser, h.GetManyNews)
	// end news

}

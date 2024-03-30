package main

import (
	"fmt"
	"github.com/JohnKucharsky/DevGroup/db"
	"github.com/JohnKucharsky/DevGroup/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// ghg
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Can't load env")
	}
	postgresURI := os.Getenv("POSTGRES_URI")
	redisURI := os.Getenv("REDIS_URI")

	f := fiber.New()
	f.Use(logger.New())
	f.Use(
		cors.New(
			cors.Config{
				AllowOrigins: "*",
				AllowHeaders: "Origin, Content-Type, Accept, Authorization",
				AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE",
			},
		),
	)
	// check
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = f.Shutdown()
	}()
	d := db.New(postgresURI)
	redis := db.NewRedis(redisURI)

	router.Register(f, d, redis)
	f.Get(
		"/api", func(c *fiber.Ctx) error {
			if err := c.SendFile("./public/index.html"); err != nil {
				err := c.SendStatus(http.StatusBadRequest)
				if err != nil {
					return nil
				}
			}
			return nil
		},
	)

	err := f.Listen(":8080")
	if err != nil {
		fmt.Printf("%v", err.Error())
	}
	fmt.Println("Running cleanup tasks...")
}

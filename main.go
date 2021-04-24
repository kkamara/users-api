package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kkamara/users-api/handlers"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		return c.Next()
	})

	api := app.Group("/api")

	api.Post("/users", handlers.PostUser)
	api.Put("/users/:username", handlers.PutUser)

	app.Listen(":3000")
}

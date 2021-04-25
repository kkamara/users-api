package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kkamara/users-api/handlers"
	"github.com/kkamara/users-api/models/userModel"
)

func main() {
	app := *fiber.New()

	app.Use(logger.New())
	app.Use(func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		return c.Next()
	})
	app.Use(func(c *fiber.Ctx) error {
		if c.OriginalURL() == "/api/users" && c.Method() == "POST" {
			return c.Next()
		}
		if authKey := c.Get("Authorization"); len(authKey) > 0 {
			_, err := userModel.VerifyAuthToken(authKey)
			if err != nil {
				fmt.Printf("%v", err)
				return c.JSON(fiber.Map{"error": "Unauthorized"})
			}
			return c.Next()
		}
		c.Context().SetStatusCode(401)
		return c.JSON(fiber.Map{"error": "Unauthorized"})
	})

	err := userModel.Seed()
	if err != nil {
		panic(err)
	}

	api := app.Group("/api")

	api.Post("/users", handlers.PostUser)
	api.Get("/users", handlers.GetUsers)
	api.Get("/users/search", handlers.SearchUsers)
	api.Patch("/users/:username", handlers.PatchUser)
	api.Put("/users/:username/darkmode", handlers.PutToggleDarkMode)
	api.Delete("/users/:username", handlers.DeleteUser)

	app.Listen(":3000")
}

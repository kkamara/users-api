package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kkamara/users-api/models/userModel"
	"github.com/kkamara/users-api/schemas/userSchema"
)

func PostUser(c *fiber.Ctx) error {
	user := new(userSchema.UserSchema)
	if err := c.BodyParser(user); err != nil {
		c.Context().SetStatusCode(500)
		return c.JSON(fiber.Map{"error": "Error encountered when parsing request input."})
	}

	user.Username = userModel.GenerateUsername(user.FirstName, user.LastName)

	const createdFormat = "2006-01-02 15:04:05"
	user.DateCreated = time.Now().Format(createdFormat)

	errors := userModel.ValidateCreate(user)
	if len(errors) > 0 {
		c.Context().SetStatusCode(400)
		return c.JSON(fiber.Map{"errors": errors})
	}
	newUser, err := userModel.Create(user)
	if err != nil {
		c.Context().SetStatusCode(500)
		return c.JSON(fiber.Map{"error": "Unknown error encountered when saving resource."})
	}
	return c.JSON(fiber.Map{"data": &newUser})
}

func PutUser(c *fiber.Ctx) error {
	user, err := userModel.Get(c.Params("username"))
	if err != nil {
		c.Context().SetStatusCode(404)
		return c.JSON(fiber.Map{"error": "Resource does not exist."})
	}
	if err := c.BodyParser(user); err != nil {
		c.Context().SetStatusCode(500)
		return c.JSON(fiber.Map{"error": "Error encountered when parsing request input."})
	}

	errors := userModel.ValidateUpdate(user)
	if len(errors) > 0 {
		c.Context().SetStatusCode(400)
		return c.JSON(fiber.Map{"errors": errors})
	}
	newUser, err := userModel.Update(user.Username, user)
	if err != nil {
		c.Context().SetStatusCode(500)
		return c.JSON(fiber.Map{"error": "Unknown error encountered when saving resource."})
	}
	return c.JSON(fiber.Map{"data": &newUser})
}

func GetUsers(c *fiber.Ctx) error {
	users, err := userModel.GetAll()
	if err != nil {
		c.Context().SetStatusCode(500)
		return c.JSON(fiber.Map{"error": "Unknown error encountered when saving resource."})
	}
	return c.JSON(fiber.Map{"data": &users})
}

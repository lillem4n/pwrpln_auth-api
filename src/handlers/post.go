package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gitlab.larvit.se/power-plan/api/src/db"
)

// UserCreate creates a new user
func UserCreate(c *fiber.Ctx) error {
	type UserInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	userInput := new(UserInput)

	if err := c.BodyParser(userInput); err != nil {
		return c.Status(400).JSON([]ResJSONError{
			{Error: err.Error()},
		})
	}

	var errors []ResJSONError

	if userInput.Username == "" {
		errors = append(errors, ResJSONError{Error: "Can not be empty", Field: "username"})
	}
	if userInput.Password == "" {
		errors = append(errors, ResJSONError{Error: "Can not be empty", Field: "password"})
	}

	if len(errors) != 0 {
		return c.Status(400).JSON(errors)
	}

	createdUser := db.User{
		ID:       uuid.New(),
		Username: userInput.Username,
	}

	return c.Status(201).JSON(createdUser)
}

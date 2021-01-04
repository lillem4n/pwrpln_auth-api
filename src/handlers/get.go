package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// Hello handler
func (h Handlers) Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

// AccountGet handler
func (h Handlers) AccountGet(c *fiber.Ctx) error {
	accountID := c.Params("accountID")
	// logContext := log.WithFields(log.Fields{"accountID": accountID})

	authErr := h.RequireAdminRoleOrAccountID(c, accountID)
	if authErr != nil {
		return c.Status(403).JSON([]ResJSONError{{Error: authErr.Error()}})
	}

	account, accountErr := h.Db.AccountGet(accountID, "")
	if accountErr != nil {
		return c.Status(500).JSON([]ResJSONError{{Error: accountErr.Error()}})
	}

	return c.JSON(account)
}

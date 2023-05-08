package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AccountDel godoc
// @Summary Delete an account
// @Description Requires Authorization-header with role "admin" or a matching account id
// @Description Example: Authorization: bearer xxx
// @Description Where "xxx" is a valid JWT token
// @ID account-del
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Success 204 {string} string ""
// @Failure 400 {object} []ResJSONError
// @Failure 401 {object} []ResJSONError
// @Failure 403 {object} []ResJSONError
// @Failure 404 {object} []ResJSONError
// @Failure 415 {object} []ResJSONError
// @Failure 500 {object} []ResJSONError
// @Router /accounts/:id [delete]
func (h Handlers) AccountDel(c *fiber.Ctx) error {
	accountID := c.Params("accountID")

	_, uuidErr := uuid.Parse(accountID)
	if uuidErr != nil {
		return c.Status(400).JSON([]ResJSONError{{Error: "Invalid uuid format"}})
	}

	authErr := h.RequireAdminRoleOrAccountID(c, accountID)
	if authErr != nil {
		return c.Status(403).JSON([]ResJSONError{{Error: authErr.Error()}})
	}

	err := h.Db.AccountDel(accountID)
	if err != nil {
		if err.Error() == "No account found for given accountID" {
			return c.Status(404).JSON([]ResJSONError{{Error: err.Error()}})
		} else {
			return c.Status(500).JSON([]ResJSONError{{Error: "Database error when trying to remove account"}})
		}
	}

	return c.Status(204).Send(nil)
}

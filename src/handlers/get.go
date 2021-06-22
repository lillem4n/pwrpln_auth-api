package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// AccountGet godoc
// @Summary Get account by id
// @Description Requires Authorization-header with either role "admin" or with a matching account id.
// @Description Example: Authorization: bearer xxx
// @Description Where "xxx" is a valid JWT token
// @ID get-account-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Success 200 {object} db.Account
// @Failure 401 {object} []ResJSONError
// @Failure 403 {object} []ResJSONError
// @Failure 415 {object} []ResJSONError
// @Failure 500 {object} []ResJSONError
// @Router /account/{id} [get]
func (h Handlers) AccountGet(c *fiber.Ctx) error {
	accountID := c.Params("accountID")
	// logContext := log.WithFields(log.Fields{"accountID": accountID})

	authErr := h.RequireAdminRoleOrAccountID(c, accountID)
	if authErr != nil {
		return c.Status(403).JSON([]ResJSONError{{Error: authErr.Error()}})
	}

	account, accountErr := h.Db.AccountGet(accountID, "", "")
	if accountErr != nil {
		return c.Status(500).JSON([]ResJSONError{{Error: accountErr.Error()}})
	}

	return c.JSON(account)
}

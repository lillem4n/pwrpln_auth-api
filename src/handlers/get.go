package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AccountGet godoc
// @Summary Get account by id
// @Description Requires Authorization-header with either role "admin" or with a matching account id.
// @Description Example: Authorization: bearer xxx
// @Description Where "xxx" is a valid JWT token
// @Param Authorization header string true "Insert your access token"
// @ID get-account-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Success 200 {object} db.Account
// @Failure 401 {object} []ResJSONError
// @Failure 403 {object} []ResJSONError
// @Failure 415 {object} []ResJSONError
// @Failure 500 {object} []ResJSONError
// @Router /accounts/{id} [get]
func (h Handlers) AccountGet(c *fiber.Ctx) error {
	accountID := c.Params("accountID")

	_, uuidErr := uuid.Parse(accountID)
	if uuidErr != nil {
		return c.Status(400).JSON([]ResJSONError{{Error: "Invalid uuid format"}})
	}

	authErr := h.RequireAdminRoleOrAccountID(c, accountID)
	if authErr != nil {
		return c.Status(403).JSON([]ResJSONError{{Error: authErr.Error()}})
	}

	account, accountErr := h.Db.AccountGet(accountID, "", "")
	if accountErr != nil {
		if accountErr.Error() == "no rows in result set" {
			return c.Status(404).JSON([]ResJSONError{{Error: "No account found for given accountID"}})
		} else {
			return c.Status(500).JSON([]ResJSONError{{Error: accountErr.Error()}})
		}
	}

	return c.JSON(account)
}

// AccountGet godoc
// @Summary Get accounts
// @Description Requires Authorization-header with role "admin".
// @Description Example: Authorization: bearer xxx
// @Description Where "xxx" is a valid JWT token
// @Param Authorization header string true "Insert your access token"
// @Accept  json
// @Produce  json
// @Success 200 {object} []db.Account
// @Failure 401 {object} []ResJSONError
// @Failure 403 {object} []ResJSONError
// @Failure 415 {object} []ResJSONError
// @Failure 500 {object} []ResJSONError
// @Router /accounts [get]
func (h Handlers) AccountsGet(c *fiber.Ctx) error {
	accountID := c.Params("accountID")

	authErr := h.RequireAdminRole(c)
	if authErr != nil {
		return c.Status(403).JSON([]ResJSONError{{Error: authErr.Error()}})
	}

	account, accountErr := h.Db.AccountGet(accountID, "", "")
	if accountErr != nil {
		if accountErr.Error() == "no rows in result set" {
			return c.Status(404).JSON([]ResJSONError{{Error: "No account found for given accountID"}})
		} else {
			return c.Status(500).JSON([]ResJSONError{{Error: accountErr.Error()}})
		}
	}

	return c.JSON(account)
}

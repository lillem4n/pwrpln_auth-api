package handlers

import (
	"gitea.larvit.se/pwrpln/auth-api/src/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AccountUpdateFields godoc
// @Summary Update account fields
// @Description Requires Authorization-header with role "admin".
// @Description Example: Authorization: bearer xxx
// @Description Where "xxx" is a valid JWT token
// @Param Authorization header string true "Insert your access token"
// @ID account-update-fields
// @Accept  json
// @Produce  json
// @Param body body []db.AccountCreateInputFields true "Fields array with objects to be written to database"
// @Success 200 {object} db.Account
// @Failure 400 {object} []ResJSONError
// @Failure 401 {object} []ResJSONError
// @Failure 403 {object} []ResJSONError
// @Failure 415 {object} []ResJSONError
// @Failure 500 {object} []ResJSONError
// @Router /accounts/{id}/fields [put]
func (h Handlers) AccountUpdateFields(c *fiber.Ctx) error {
	accountID := c.Params("accountID")

	h.Log.Context = []interface{}{
		"accountID", accountID,
	}

	_, uuidErr := uuid.Parse(accountID)
	if uuidErr != nil {
		h.Log.Debug("client supplied invalid uuid format")
		return c.Status(400).JSON([]ResJSONError{{Error: "Invalid uuid format"}})
	}

	authErr := h.RequireAdminRole(c)
	if authErr != nil {
		h.Log.Debug("client does not have admin role")
		return c.Status(403).JSON([]ResJSONError{{Error: authErr.Error()}})
	}

	fieldsInput := new([]db.AccountCreateInputFields)

	if err := c.BodyParser(fieldsInput); err != nil {
		return c.Status(400).JSON([]ResJSONError{
			{Error: err.Error()},
		})
	}

	updatedAccount, err := h.Db.AccountUpdateFields(accountID, *fieldsInput)
	if err != nil {
		return c.Status(500).JSON([]ResJSONError{{Error: "Internal server error"}})
	}

	return c.Status(200).JSON(updatedAccount)
}

package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.larvit.se/power-plan/auth/src/db"
	"gitlab.larvit.se/power-plan/auth/src/utils"
)

// AccountCreate godoc
// @Summary Create an account
// @Description Create an account
// @ID account-create
// @Accept  json
// @Produce  json
// @Success 200 {object} db.CreatedAccount
// @Failure 401 {object} ResJSONError
// @Failure 403 {object} ResJSONError
// @Failure 415 {object} ResJSONError
// @Failure 500 {object} ResJSONError
// @Router /account [post]
func (h Handlers) AccountCreate(c *fiber.Ctx) error {
	authErr := h.RequireAdminRole(c)
	if authErr != nil {
		return c.Status(403).JSON([]ResJSONError{{Error: authErr.Error()}})
	}

	type AccountInput struct {
		Name     string                        `json:"name"`
		Password string                        `json:"password"`
		Fields   []db.AccountCreateInputFields `json:"fields"`
	}

	accountInput := new(AccountInput)

	if err := c.BodyParser(accountInput); err != nil {
		return c.Status(400).JSON([]ResJSONError{
			{Error: err.Error()},
		})
	}

	var errors []ResJSONError

	if accountInput.Name == "" {
		errors = append(errors, ResJSONError{Error: "Can not be empty", Field: "name"})
	}

	if len(errors) != 0 {
		return c.Status(400).JSON(errors)
	}

	newAccountID, uuidErr := uuid.NewRandom()
	if uuidErr != nil {
		h.Log.Fatal("Could not create new Uuid", "err", uuidErr.Error())
	}

	hashedPwd, pwdErr := utils.HashPassword(accountInput.Password)
	if pwdErr != nil {
		h.Log.Fatal("Could not hash password", "err", pwdErr.Error())
	}

	createdAccount, err := h.Db.AccountCreate(db.AccountCreateInput{
		ID:       newAccountID,
		Name:     accountInput.Name,
		APIKey:   utils.RandString(60),
		Fields:   accountInput.Fields,
		Password: hashedPwd,
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
			return c.Status(409).JSON([]ResJSONError{{Error: "accountName is already taken"}})
		}
		return c.Status(500).JSON([]ResJSONError{{Error: err.Error()}})
	}

	return c.Status(201).JSON(createdAccount)
}

// AccountAuthAPIKey godoc
// @Summary Authenticate account by API Key
// @Description Authenticate account by API Key
// @ID auth-account-by-api-key
// @Accept  json
// @Produce  json
// @Success 200 {object} db.Account
// @Failure 401 {object} ResJSONError
// @Failure 403 {object} ResJSONError
// @Failure 415 {object} ResJSONError
// @Failure 500 {object} ResJSONError
// @Router /auth/api-key [post]
func (h Handlers) AccountAuthAPIKey(c *fiber.Ctx) error {
	inputAPIKey := string(c.Request().Body())
	inputAPIKey = inputAPIKey[1 : len(inputAPIKey)-1]

	resolvedAccount, accountErr := h.Db.AccountGet("", inputAPIKey, "")
	if accountErr != nil {
		if accountErr.Error() == "no rows in result set" {
			return c.Status(403).JSON([]ResJSONError{{Error: "Invalid credentials"}})
		}
		h.Log.Error("Something went wrong when trying to fetch account", "err", accountErr.Error())
		return c.Status(500).JSON([]ResJSONError{{Error: "Something went wrong when trying to fetch account"}})
	}

	return h.returnTokens(resolvedAccount, c)
}

// AccountAuthPassword godoc
// @Summary Authenticate account by Password
// @Description Authenticate account by Password
// @ID auth-account-by-password
// @Accept  json
// @Produce  json
// @Success 200 {object} db.Account
// @Failure 401 {object} ResJSONError
// @Failure 403 {object} ResJSONError
// @Failure 415 {object} ResJSONError
// @Failure 500 {object} ResJSONError
// @Router /auth/password [post]
func (h Handlers) AccountAuthPassword(c *fiber.Ctx) error {
	type AuthInput struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	authInput := new(AuthInput)
	if err := c.BodyParser(authInput); err != nil {
		return c.Status(400).JSON([]ResJSONError{{Error: err.Error()}})
	}

	resolvedAccount, err := h.Db.AccountGet("", "", authInput.Name)
	if err != nil {
		if err.Error() == "No account found" {
			return c.Status(403).JSON([]ResJSONError{{Error: "Invalid name or password"}})
		}

		return c.Status(500).JSON([]ResJSONError{{Error: err.Error()}})
	}

	if utils.CheckPasswordHash(authInput.Password, resolvedAccount.Password) == false {
		return c.Status(403).JSON([]ResJSONError{{Error: "Invalid name or password"}})
	}

	return h.returnTokens(resolvedAccount, c)
}

// RenewToken godoc
// @Summary Renew token
// @Description Renew token
// @ID renew-token
// @Accept  json
// @Produce  json
// @Success 200 {object} db.Account
// @Failure 401 {object} ResJSONError
// @Failure 403 {object} ResJSONError
// @Failure 415 {object} ResJSONError
// @Failure 500 {object} ResJSONError
// @Router /renew-token [post]
func (h Handlers) RenewToken(c *fiber.Ctx) error {
	inputToken := string(c.Request().Body())
	inputToken = inputToken[1 : len(inputToken)-1]

	foundAccountID, err := h.Db.RenewalTokenGet(inputToken)
	if err != nil {
		return c.Status(500).JSON([]ResJSONError{{Error: err.Error()}})
	} else if foundAccountID == "" {
		return c.Status(403).JSON([]ResJSONError{{Error: "Invalid token"}})
	}

	resolvedAccount, accountErr := h.Db.AccountGet(foundAccountID, "", "")
	if accountErr != nil {
		if accountErr.Error() == "no rows in result set" {
			return c.Status(500).JSON([]ResJSONError{{Error: "Database missmatch. Token found, but account is missing."}})
		}
		h.Log.Error("Something went wrong when trying to fetch account", "err", accountErr.Error())
		return c.Status(500).JSON([]ResJSONError{{Error: "Something went wrong when trying to fetch account"}})
	}

	rmErr := h.Db.RenewalTokenRm(inputToken)
	if rmErr != nil {
		h.Log.Error("Something went wrong when trying to fetch account", "err", rmErr.Error())
		return c.Status(500).JSON([]ResJSONError{{Error: "Could not remove old token"}})
	}

	return h.returnTokens(resolvedAccount, c)
}

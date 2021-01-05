package handlers

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.larvit.se/power-plan/auth/src/db"
	"gitlab.larvit.se/power-plan/auth/src/utils"
)

// AccountCreate creates a new account
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
		log.Fatal("Could not create new Uuid, err: " + uuidErr.Error())
	}

	hashedPwd, pwdErr := utils.HashPassword(accountInput.Password)
	if pwdErr != nil {
		log.Fatal("Could not hash password, err: " + pwdErr.Error())
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

// AccountAuthAPIKey auths an APIKey
func (h Handlers) AccountAuthAPIKey(c *fiber.Ctx) error {
	inputAPIKey := string(c.Request().Body())
	inputAPIKey = inputAPIKey[1 : len(inputAPIKey)-1]

	resolvedAccount, accountErr := h.Db.AccountGet("", inputAPIKey, "")
	if accountErr != nil {
		if accountErr.Error() == "no rows in result set" {
			return c.Status(403).JSON([]ResJSONError{{Error: "Invalid credentials"}})
		}
		log.Error("Something went wrong when trying to fetch account")
		return c.Status(500).JSON([]ResJSONError{{Error: "Something went wrong when trying to fetch account"}})
	}

	return h.returnTokens(resolvedAccount, c)
}

// AccountAuthPassword auths a name/password pair
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

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
	type AccountInput struct {
		AccountName string                        `json:"accountName"`
		Password    string                        `json:"password"`
		Fields      []db.AccountCreateInputFields `json:"fields"`
	}

	accountInput := new(AccountInput)

	if err := c.BodyParser(accountInput); err != nil {
		return c.Status(400).JSON([]ResJSONError{
			{Error: err.Error()},
		})
	}

	var errors []ResJSONError

	if accountInput.AccountName == "" {
		errors = append(errors, ResJSONError{Error: "Can not be empty", Field: "accountName"})
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

	apiKey := utils.RandString(60)

	createdAccount, err := h.Db.AccountCreate(db.AccountCreateInput{
		ID:          newAccountID,
		AccountName: accountInput.AccountName,
		APIKey:      apiKey,
		Fields:      accountInput.Fields,
		Password:    hashedPwd,
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
	return c.Status(200).JSON("key höhö")
}

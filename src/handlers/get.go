package handlers

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

// Hello handler
func (h Handlers) Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

// AccountGet handler
func (h Handlers) AccountGet(c *fiber.Ctx) error {
	log.WithFields(log.Fields{"accountID": c.Params("accountID")}).Debug("GETing account")
	return c.SendString("Account ffs")
}

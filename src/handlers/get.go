package handlers

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

// Hello handler
func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

// UserGet handler
func UserGet(c *fiber.Ctx) error {
	log.WithFields(log.Fields{"userID": c.Params("userID")}).Debug("GETing user")
	return c.SendString("USER ffs")
}

package handlers

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

// Log all requests
func (h Handlers) Log(c *fiber.Ctx) error {
	log.WithFields(log.Fields{
		"method": c.Method(),
		"url":    c.OriginalURL(),
	}).Debug("http request")

	c.Next()
	return nil
}

// RequireJSON is a middleware that makes sure the request content-type always is application/json (or nothing, defaulting to application/json)
func (h Handlers) RequireJSON(c *fiber.Ctx) error {
	c.Accepts("application/json")
	contentType := string(c.Request().Header.ContentType())

	if contentType != "application/json" && contentType != "" {
		log.WithFields(log.Fields{"content-type": contentType}).Debug("Invalid content-type in request")
		return c.Status(415).JSON([]ResJSONError{{Error: "Invalid content-type"}})
	}

	c.Next()
	return nil
}

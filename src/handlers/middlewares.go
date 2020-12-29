package handlers

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

// RequireJSON is a middleware that makes sure the request content-type always is application/json (or nothing, defaulting to application/json)
func RequireJSON(c *fiber.Ctx) error {
	c.Accepts("application/json")
	contentType := string(c.Request().Header.ContentType())

	if contentType != "application/json" && contentType != "" {
		log.WithFields(log.Fields{"content-type": contentType}).Debug("Invalid content-type in request")
		return c.Status(415).JSON([]ResJSONError{{Error: "Invalid content-type"}})
	}

	c.Next()
	return nil
}

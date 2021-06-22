package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// Log all requests
func (h Handlers) LogReq(c *fiber.Ctx) error {
	h.Log.Debug("http request", "method", c.Method(), "url", c.OriginalURL())

	c.Next()
	return nil
}

// RequireJSON is a middleware that makes sure the request content-type always is application/json (or nothing, defaulting to application/json)
func (h Handlers) RequireJSON(c *fiber.Ctx) error {
	c.Accepts("application/json")
	contentType := string(c.Request().Header.ContentType())

	if contentType != "application/json" && contentType != "" {
		h.Log.Debug("Invalid content-type in request", "content-type", contentType)
		return c.Status(415).JSON([]ResJSONError{{Error: "Invalid content-type"}})
	}

	c.Next()
	return nil
}

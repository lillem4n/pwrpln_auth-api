package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"gitlab.larvit.se/power-plan/api/src/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Add this line for logging filename and line number!
	// log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	app := fiber.New()

	// Always require application/json
	app.Use(handlers.RequireJSON)

	app.Get("/", handlers.Hello)
	app.Get("/user/:userID", handlers.UserGet)
	app.Post("/user", handlers.UserCreate)

	log.WithFields(log.Fields{"WEB_BIND_HOST": os.Getenv("WEB_BIND_HOST")}).Debug("Trying to start web server")

	app.Listen(os.Getenv("WEB_BIND_HOST"))
}

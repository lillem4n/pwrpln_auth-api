package main

import (
	"context"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gitlab.larvit.se/power-plan/auth/src/db"
	h "gitlab.larvit.se/power-plan/auth/src/handlers"
)

func createAdminAccount(Db db.Db) {
	adminAccountID, uuidErr := uuid.NewRandom()
	if uuidErr != nil {
		log.Fatal("Could not create new Uuid, err: " + uuidErr.Error())
	}
	_, adminAccountErr := Db.AccountCreate(db.AccountCreateInput{
		ID:       adminAccountID,
		Name:     "admin",
		APIKey:   os.Getenv("ADMIN_API_KEY"),
		Password: "",
		Fields:   []db.AccountCreateInputFields{{Name: "role", Values: []string{"admin"}}},
	})
	if adminAccountErr != nil && strings.HasPrefix(adminAccountErr.Error(), "ERROR: duplicate key") {
		log.Info("Admin account already created, nothing written to database")
	} else if adminAccountErr != nil {
		log.Fatal("Could not create admin account, err: " + adminAccountErr.Error())
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Add this line for logging filename and line number!
	// log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	if os.Getenv("JWT_SHARED_SECRET") == "changeMe" {
		log.Fatal("You must change JWT_SHARED_SECRET in .env")
	}
	if os.Getenv("ADMIN_API_KEY") == "changeMe" {
		log.Fatal("You must change ADMIN_API_KEY in .env")
	}
	jwtKey := []byte(os.Getenv("JWT_SHARED_SECRET"))

	dbPool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to open DB connection: ", err)
	} else {
		log.Info("Connected to PostgreSQL database")
	}
	defer dbPool.Close()

	app := fiber.New()

	Db := db.Db{DbPool: dbPool}
	handlers := h.Handlers{Db: Db, JwtKey: jwtKey}

	createAdminAccount(Db)

	// Log all requests
	app.Use(handlers.Log)

	// Always require application/json
	app.Use(handlers.RequireJSON)

	app.Get("/", handlers.Hello)
	app.Get("/account/:accountID", handlers.AccountGet)
	app.Post("/account", handlers.AccountCreate)
	app.Post("/auth/api-key", handlers.AccountAuthAPIKey)
	app.Post("/auth/password", handlers.AccountAuthPassword)
	app.Post("/renew-token", handlers.TokenRenew)

	log.WithFields(log.Fields{"WEB_BIND_HOST": os.Getenv("WEB_BIND_HOST")}).Info("Trying to start web server")

	if err := app.Listen(os.Getenv("WEB_BIND_HOST")); err != nil {
		log.Fatal(err)
	}

	log.Info("Webb server closed, shutting down")
}

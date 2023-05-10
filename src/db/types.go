package db

import (
	"time"

	"gitea.larvit.se/pwrpln/go_log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Account is an account as represented in the database
type Account struct {
	ID       uuid.UUID           `json:"id"`
	Created  time.Time           `json:"created"`
	Fields   map[string][]string `json:"fields"`
	Name     string              `json:"name"`
	Password string              `json:"-"`
}

// CreatedAccount is a newly created account in the system
type CreatedAccount struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	APIKey string    `json:"apiKey"`
}

// AccountCreateInputFields yes
type AccountCreateInputFields struct {
	Name   string
	Values []string
}

// AccountCreateInput is used as input struct for database creation of account
type AccountCreateInput struct {
	ID       uuid.UUID
	Name     string
	APIKey   string
	Fields   []AccountCreateInputFields
	Password string
}

// Db struct
type Db struct {
	DbPool *pgxpool.Pool
	Log    go_log.Log
}

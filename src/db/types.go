package db

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// CreatedAccount is a newly created account in the system
type CreatedAccount struct {
	ID          uuid.UUID `json:"id"`
	AccountName string    `json:"accountName"`
	APIKey      string    `json:"apiKey"`
}

// AccountCreateInputFields yes
type AccountCreateInputFields struct {
	Name   string
	Values []string
}

// AccountCreateInput is used as input struct for database creation of account
type AccountCreateInput struct {
	ID          uuid.UUID
	AccountName string
	APIKey      string
	Fields      []AccountCreateInputFields
	Password    string
}

// Db struct
type Db struct {
	DbPool *pgxpool.Pool
}

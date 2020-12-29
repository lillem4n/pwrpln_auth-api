package db

import (
	"github.com/google/uuid"
)

// User is a user in the system
type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

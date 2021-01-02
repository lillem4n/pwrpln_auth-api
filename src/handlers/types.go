package handlers

import (
	"gitlab.larvit.se/power-plan/auth/src/db"
)

// Handlers is the overall struct for all http request handlers
type Handlers struct {
	Db db.Db
}

// ResJSONError is an error field that is used in JSON error responses
type ResJSONError struct {
	Error string `json:"error"`
	Field string `json:"field,omitempty"`
}

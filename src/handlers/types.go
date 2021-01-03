package handlers

import (
	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.larvit.se/power-plan/auth/src/db"
)

// Claims is the JWT struct
type Claims struct {
	AccountID     string              `json:"accountId"`
	AccountFields map[string][]string `json:"accountFields"`
	AccountName   string              `json:"accountName"`
	jwt.StandardClaims
}

// Handlers is the overall struct for all http request handlers
type Handlers struct {
	Db     db.Db
	JwtKey []byte
}

// ResJSONError is an error field that is used in JSON error responses
type ResJSONError struct {
	Error string `json:"error"`
	Field string `json:"field,omitempty"`
}

// ResToken is a response used to return a valid token and valid renewalToken
type ResToken struct {
	JWT          string `json:"jwt"`
	RenewalToken string `json:"renewalToken"`
}

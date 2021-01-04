package handlers

import (
	"errors"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"gitlab.larvit.se/power-plan/auth/src/db"
)

func (h Handlers) returnTokens(account db.Account, c *fiber.Ctx) error {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		AccountID:     account.ID.String(),
		AccountName:   account.Name,
		AccountFields: account.Fields,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.JwtKey)
	if err != nil {
		log.Error("Could not create token string, err: " + err.Error())
		return c.Status(500).JSON([]ResJSONError{{Error: "Could not create JWT token string"}})
	}

	renewalToken, renewalTokenErr := h.Db.RenewalTokenGet(account.ID.String())
	if renewalTokenErr != nil {
		log.Error("Could not create renewal token, err: " + renewalTokenErr.Error())
		return c.Status(500).JSON([]ResJSONError{{Error: "Could not create renewal token"}})
	}

	return c.Status(200).JSON(ResToken{
		JWT:          tokenString,
		RenewalToken: renewalToken,
	})
}

func (h Handlers) parseJWT(JWT string) (Claims, error) {
	logContext := log.WithFields(log.Fields{"JWT": JWT})
	logContext.Trace("Parsing JWT")

	JWT = strings.TrimPrefix(JWT, "bearer ") // Since the Authorization header should always start with "bearer "
	logContext.WithFields(log.Fields{"TrimmedJWT": JWT}).Trace("JWT trimmed")

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(JWT, claims, func(token *jwt.Token) (interface{}, error) {
		return h.JwtKey, nil
	})
	if err != nil {
		return Claims{}, err
	}
	if !token.Valid {
		err := errors.New("Invalid token")
		return Claims{}, err
	}

	return *claims, nil
}

func (h Handlers) parseHeaders(c *fiber.Ctx) map[string]string {
	headersMap := make(map[string]string)

	headersString := c.Request().Header.String()
	headersLines := strings.Split(headersString, "\r\n")

	for _, line := range headersLines {
		lineParts := strings.Split(line, ": ")

		if len(lineParts) == 1 {
			log.WithFields(log.Fields{"line": line}).Trace("Ignoring header line")
		} else {
			headersMap[lineParts[0]] = lineParts[1]
		}
	}

	return headersMap
}

// RequireAdminRole returns nil if no error is found
func (h Handlers) RequireAdminRole(c *fiber.Ctx) error {
	headers := h.parseHeaders(c)

	if headers["Authorization"] == "" {
		return errors.New("Authorization header is missing")
	}

	claims, claimsErr := h.parseJWT(headers["Authorization"])
	if claimsErr != nil {
		return claimsErr
	}

	if claims.AccountFields == nil {
		return errors.New("Account have no fields at all")
	}

	if claims.AccountFields["role"] == nil {
		return errors.New("Account have no field named \"role\"")
	}

	for _, role := range claims.AccountFields["role"] {
		if role == "admin" {
			return nil
		}
	}

	return errors.New("No \"admin\" role found on account")
}

// RequireAdminRoleOrAccountID returns nil if no error is found
func (h Handlers) RequireAdminRoleOrAccountID(c *fiber.Ctx, accountID string) error {
	headers := h.parseHeaders(c)

	if headers["Authorization"] == "" {
		return errors.New("Authorization header is missing")
	}

	claims, claimsErr := h.parseJWT(headers["Authorization"])
	if claimsErr != nil {
		return claimsErr
	}

	if claims.AccountID == accountID {
		return nil
	}

	if claims.AccountFields == nil {
		return errors.New("AccountID does not match and account have no fields at all")
	}

	if claims.AccountFields["role"] == nil {
		return errors.New("AccountID does not match and account have no field named \"role\"")
	}

	for _, role := range claims.AccountFields["role"] {
		if role == "admin" {
			return nil
		}
	}

	return errors.New("AccountID does not match and no \"admin\" role found on account")
}

package db

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.larvit.se/power-plan/auth/src/utils"
)

// RenewalTokenCreate obtain a new renewal token
func (d Db) RenewalTokenCreate(accountID string) (string, error) {
	logContext := log.WithFields(log.Fields{"accountID": accountID})

	logContext.Debug("Creating new renewal token")

	newToken := utils.RandString(60)

	insertSQL := "INSERT INTO \"renewalTokens\" (\"accountId\",token) VALUES($1,$2);"
	_, insertErr := d.DbPool.Exec(context.Background(), insertSQL, accountID, newToken)
	if insertErr != nil {
		logContext.Error("Could not insert into database table \"renewalTokens\", err: " + insertErr.Error())
		return "", insertErr
	}

	return newToken, nil
}

// RenewalTokenGet checks if a valid renewal token exists in database
func (d Db) RenewalTokenGet(token string) (string, error) {
	log.Debug("Trying to get a renewal token")

	sql := "SELECT \"accountId\" FROM \"renewalTokens\" WHERE exp >= now() AND token = $1"

	var foundAccountID string
	err := d.DbPool.QueryRow(context.Background(), sql, token).Scan(&foundAccountID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", nil
		}

		log.Error("Database error when fetching renewal token, err: " + err.Error())
		return "", err
	}

	return foundAccountID, nil
}

// RenewalTokenRm removes a renewal token from the database
func (d Db) RenewalTokenRm(token string) error {
	log.Debug("Trying to remove a renewal token")

	sql := "DELETE FROM \"renewalTokens\" WHERE token = $1"
	_, err := d.DbPool.Exec(context.Background(), sql, token)
	if err != nil {
		log.Error("Database error when trying to remove token, err: " + err.Error())
		return err
	}

	return nil
}

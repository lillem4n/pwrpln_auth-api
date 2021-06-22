package db

import (
	"context"

	"gitlab.larvit.se/power-plan/auth/src/utils"
)

// RenewalTokenCreate obtain a new renewal token
func (d Db) RenewalTokenCreate(accountID string) (string, error) {
	d.Log.Debug("Creating new renewal token", "accountID", accountID)

	newToken := utils.RandString(60)

	insertSQL := "INSERT INTO \"renewalTokens\" (\"accountId\",token) VALUES($1,$2);"
	_, insertErr := d.DbPool.Exec(context.Background(), insertSQL, accountID, newToken)
	if insertErr != nil {
		d.Log.Error("Could not insert into database table \"renewalTokens\"", "err", insertErr.Error(), "accountID", accountID)
		return "", insertErr
	}

	return newToken, nil
}

// RenewalTokenGet checks if a valid renewal token exists in database
func (d Db) RenewalTokenGet(token string) (string, error) {
	d.Log.Debug("Trying to get a renewal token")

	sql := "SELECT \"accountId\" FROM \"renewalTokens\" WHERE exp >= now() AND token = $1"

	var foundAccountID string
	err := d.DbPool.QueryRow(context.Background(), sql, token).Scan(&foundAccountID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", nil
		}

		d.Log.Error("Database error when fetching renewal token", "err", err.Error())
		return "", err
	}

	return foundAccountID, nil
}

// RenewalTokenRm removes a renewal token from the database
func (d Db) RenewalTokenRm(token string) error {
	d.Log.Debug("Trying to remove a renewal token")

	sql := "DELETE FROM \"renewalTokens\" WHERE token = $1"
	_, err := d.DbPool.Exec(context.Background(), sql, token)
	if err != nil {
		d.Log.Error("Database error when trying to remove token", "err", err.Error())
		return err
	}

	return nil
}

package db

import (
	"context"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// AccountCreate writes a user to database
func (d Db) AccountCreate(input AccountCreateInput) (CreatedAccount, error) {
	accountSQL := "INSERT INTO accounts (id, name, \"apiKey\", password) VALUES($1,$2,$3,$4);"

	_, err := d.DbPool.Exec(context.Background(), accountSQL, input.ID, input.Name, input.APIKey, input.Password)
	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
			log.WithFields(log.Fields{"name": input.Name}).Debug("Duplicate name in accounts database")
		} else {
			log.Error("Database error when trying to add account: " + err.Error())
		}

		return CreatedAccount{}, err
	}

	log.WithFields(log.Fields{
		"id":   input.ID,
		"name": input.Name,
	}).Info("Added account to database")

	accountFieldsSQL := "INSERT INTO \"accountsFields\" (id, \"accountId\", name, value) VALUES($1,$2,$3,$4);"
	for _, field := range input.Fields {
		newFieldID, uuidErr := uuid.NewRandom()
		if uuidErr != nil {
			log.Fatal("Could not create new Uuid, err: " + uuidErr.Error())
		}

		_, err := d.DbPool.Exec(context.Background(), accountFieldsSQL, newFieldID, input.ID, field.Name, field.Values)
		if err != nil {
			if strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
				log.Error("Database error when trying to account field: " + err.Error())
			}
		}

		log.WithFields(log.Fields{
			"accountId":   input.ID,
			"fieldName":   field.Name,
			"fieldValues": field.Values,
		}).Debug("Added account field")
	}

	return CreatedAccount{
		ID:     input.ID,
		Name:   input.Name,
		APIKey: input.APIKey,
	}, nil
}

// AccountGet fetches an account from the database
func (d Db) AccountGet(accountID string, APIKey string, Name string) (Account, error) {
	logContext := log.WithFields(log.Fields{
		"accountID": accountID,
		"APIKey":    len(APIKey),
	})

	logContext.Debug("Trying to get account")

	var account Account
	var searchParam string
	accountSQL := "SELECT id, created, name, \"password\" FROM accounts WHERE "
	if accountID != "" {
		accountSQL = accountSQL + "id = $1"
		searchParam = accountID
	} else if APIKey != "" {
		accountSQL = accountSQL + "\"apiKey\" = $1"
		searchParam = APIKey
	} else if Name != "" {
		accountSQL = accountSQL + "name = $1"
		searchParam = Name
	}

	accountErr := d.DbPool.QueryRow(context.Background(), accountSQL, searchParam).Scan(&account.ID, &account.Created, &account.Name, &account.Password)
	if accountErr != nil {
		if accountErr.Error() == "no rows in result set" {
			logContext.Debug("No account found")
			return Account{}, accountErr
		}

		logContext.Error("Database error when fetching account, err: " + accountErr.Error())
		return Account{}, accountErr
	}

	fieldsSQL := "SELECT name, value FROM \"accountsFields\" WHERE \"accountId\" = $1"
	rows, fieldsErr := d.DbPool.Query(context.Background(), fieldsSQL, account.ID)
	if fieldsErr != nil {
		logContext.Error("Database error when fetching account fields, err: " + accountErr.Error())
		return Account{}, fieldsErr
	}

	account.Fields = make(map[string][]string)
	for rows.Next() {
		var name string
		var value []string
		err := rows.Scan(&name, &value)
		if err != nil {
			logContext.Error("Could not get name or value from database row, err: " + err.Error())
			return Account{}, err
		}
		account.Fields[name] = value
	}

	return account, nil
}

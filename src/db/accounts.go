package db

import (
	"context"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// AccountCreate writes a user to database
func (d Db) AccountCreate(input AccountCreateInput) (CreatedAccount, error) {
	accountSQL := "INSERT INTO accounts (id, \"accountName\", \"apiKey\", password) VALUES($1,$2,$3,$4);"

	_, err := d.DbPool.Exec(context.Background(), accountSQL, input.ID, input.AccountName, input.APIKey, input.Password)
	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
			log.WithFields(log.Fields{"accountName": input.AccountName}).Debug("Duplicate accountName in accounts database")
		} else {
			log.Error("Database error when trying to add account: " + err.Error())
		}

		return CreatedAccount{}, err
	}

	log.WithFields(log.Fields{
		"id":          input.ID,
		"accountName": input.AccountName,
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
		ID:          input.ID,
		AccountName: input.AccountName,
		APIKey:      input.APIKey,
	}, nil
}

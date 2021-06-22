package db

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

// AccountCreate writes a user to database
func (d Db) AccountCreate(input AccountCreateInput) (CreatedAccount, error) {
	accountSQL := "INSERT INTO accounts (id, name, \"apiKey\", password) VALUES($1,$2,$3,$4);"

	_, err := d.DbPool.Exec(context.Background(), accountSQL, input.ID, input.Name, input.APIKey, input.Password)
	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
			d.Log.Debug("Duplicate name in accounts database", "name", input.Name)
		} else {
			d.Log.Error("Database error when trying to add account", "err", err.Error())
		}

		return CreatedAccount{}, err
	}

	d.Log.Info("Added account to database", "id", input.ID, "name", input.Name)

	accountFieldsSQL := "INSERT INTO \"accountsFields\" (id, \"accountId\", name, value) VALUES($1,$2,$3,$4);"
	for _, field := range input.Fields {
		newFieldID, uuidErr := uuid.NewRandom()
		if uuidErr != nil {
			d.Log.Fatal("Could not create new Uuid", "err", uuidErr.Error())
		}

		_, err := d.DbPool.Exec(context.Background(), accountFieldsSQL, newFieldID, input.ID, field.Name, field.Values)
		if err != nil {
			if strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
				d.Log.Error("Database error when trying to account field", "err", err.Error())
			}
		}

		d.Log.Debug("Added account field", "accountId", input.ID, "fieldName", field.Name, "fieldValues", field.Values)
	}

	return CreatedAccount{
		ID:     input.ID,
		Name:   input.Name,
		APIKey: input.APIKey,
	}, nil
}

// AccountGet fetches an account from the database
func (d Db) AccountGet(accountID string, APIKey string, Name string) (Account, error) {
	d.Log.Debug("Trying to get account", "accountID", accountID, "len(APIKey)", len(APIKey))

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
			d.Log.Debug("No account found", "accountID", accountID, "APIKey", len(APIKey))
			return Account{}, accountErr
		}

		d.Log.Error("Database error when fetching account", "err", accountErr.Error(), "accountID", accountID, "APIKey", len(APIKey))
		return Account{}, accountErr
	}

	fieldsSQL := "SELECT name, value FROM \"accountsFields\" WHERE \"accountId\" = $1"
	rows, fieldsErr := d.DbPool.Query(context.Background(), fieldsSQL, account.ID)
	if fieldsErr != nil {
		d.Log.Error("Database error when fetching account fields", "err", accountErr.Error(), "accountID", accountID, "APIKey", len(APIKey))
		return Account{}, fieldsErr
	}

	account.Fields = make(map[string][]string)
	for rows.Next() {
		var name string
		var value []string
		err := rows.Scan(&name, &value)
		if err != nil {
			d.Log.Error("Could not get name or value from database row", "err", err.Error(), "accountID", accountID, "APIKey", len(APIKey))
			return Account{}, err
		}
		account.Fields[name] = value
	}

	return account, nil
}

package db

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AccountCreate writes a user to database
func (d Db) AccountCreate(input AccountCreateInput) (CreatedAccount, error) {
	d.Log.Context = []interface{}{
		"accountName", input.Name,
		"id", input.ID,
	}
	accountSQL := "INSERT INTO accounts (id, name, \"apiKey\", password) VALUES($1,$2,$3,$4);"

	_, err := d.DbPool.Exec(context.Background(), accountSQL, input.ID, input.Name, input.APIKey, input.Password)
	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
			d.Log.Debug("Duplicate name in accounts database")
		} else {
			d.Log.Warn("Database error when trying to add account", "err", err.Error())
		}

		return CreatedAccount{}, err
	}

	d.Log.Verbose("Added account to database", "id", input.ID)

	accountFieldsSQL := "INSERT INTO \"accountsFields\" (id, \"accountId\", name, value) VALUES($1,$2,$3,$4);"
	for _, field := range input.Fields {
		newFieldID, uuidErr := uuid.NewRandom()
		if uuidErr != nil {
			d.Log.Warn("Could not create new Uuid", "err", uuidErr.Error())
			return CreatedAccount{}, uuidErr
		}

		_, err := d.DbPool.Exec(context.Background(), accountFieldsSQL, newFieldID, input.ID, field.Name, field.Values)
		if err != nil {
			//if strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
			d.Log.Warn("Database error when trying to add account field", "err", err.Error(), "accountID", input.ID, "fieldName", field.Name, "fieldvalues", field.Values)
			// }
		}

		d.Log.Debug("Added account field", "accountID", input.ID, "fieldName", field.Name, "fieldValues", field.Values)
	}

	return CreatedAccount{
		ID:     input.ID,
		Name:   input.Name,
		APIKey: input.APIKey,
	}, nil
}

func (d Db) AccountDel(accountID string) error {
	d.Log.Context = []interface{}{
		"accountID", accountID,
	}
	d.Log.Verbose("Trying to delete account")

	_, renewalTokensErr := d.DbPool.Exec(context.Background(), "DELETE FROM \"renewalTokens\" WHERE \"accountId\" = $1;", accountID)
	if renewalTokensErr != nil {
		d.Log.Error("Could not remove renewal tokens for account", "err", renewalTokensErr.Error())
		return renewalTokensErr
	}

	_, fieldsErr := d.DbPool.Exec(context.Background(), "DELETE FROM \"accountsFields\" WHERE \"accountId\" = $1;", accountID)
	if fieldsErr != nil {
		d.Log.Error("Could not remove account fields", "err", fieldsErr.Error())
		return fieldsErr
	}

	res, err := d.DbPool.Exec(context.Background(), "DELETE FROM accounts WHERE id = $1", accountID)
	if err != nil {
		d.Log.Error("Could not remove account", "err", err.Error())
		return err
	}

	if res.String() == "DELETE 0" {
		d.Log.Debug("Tried to delete account, but none exists")
		err := errors.New("no account found for given accountID")
		return err
	}

	return nil
}

// AccountGet fetches an account from the database
func (d Db) AccountGet(accountID string, APIKey string, name string) (Account, error) {
	d.Log.Context = []interface{}{
		"accountID", accountID,
		"len(APIKey)", len(APIKey),
		"name", name,
	}
	d.Log.Debug("Trying to get account")

	var account Account
	var searchParam string
	accountSQL := "SELECT id, created, name, \"password\" FROM accounts WHERE "
	if accountID != "" {
		accountSQL = accountSQL + "id = $1"
		searchParam = accountID
	} else if APIKey != "" {
		accountSQL = accountSQL + "\"apiKey\" = $1"
		searchParam = APIKey
	} else if name != "" {
		accountSQL = accountSQL + "name = $1"
		searchParam = name
	} else {
		d.Log.Debug("No get criteria entered, returning empty response without calling the database")

		return Account{}, errors.New("no rows in result set")
	}

	accountErr := d.DbPool.QueryRow(context.Background(), accountSQL, searchParam).Scan(&account.ID, &account.Created, &account.Name, &account.Password)
	if accountErr != nil {
		if accountErr.Error() == "no rows in result set" {
			d.Log.Debug("No account found")
			return Account{}, accountErr
		}

		d.Log.Error("Database error when fetching account", "err", accountErr.Error())
		return Account{}, accountErr
	}

	fieldsSQL := "SELECT name, value FROM \"accountsFields\" WHERE \"accountId\" = $1"
	rows, fieldsErr := d.DbPool.Query(context.Background(), fieldsSQL, account.ID)
	if fieldsErr != nil {
		d.Log.Error("Database error when fetching account fields", "err", accountErr.Error())
		return Account{}, fieldsErr
	}

	account.Fields = make(map[string][]string)
	for rows.Next() {
		var name string
		var value []string
		err := rows.Scan(&name, &value)
		if err != nil {
			d.Log.Error("Could not get name or value from database row", "err", err.Error())
			return Account{}, err
		}
		account.Fields[name] = value
	}

	return account, nil
}

func (d Db) AccountsGet() ([]Account, error) {
	d.Log.Debug("Trying to get accounts")

	var accountFields = make(map[uuid.UUID]map[string][]string)
	var accountIds = []string{}
	var accountsToReturn = []Account{}

	accountSQL := "SELECT id, created, name FROM accounts"

	rows, err := d.DbPool.Query(context.Background(), accountSQL)
	if err != nil {
		d.Log.Error("Error executing SELECT query for accounts", "err", err.Error())
		return []Account{}, err
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			d.Log.Error("Error while iterating accounts dataset", "err", err.Error())
			return []Account{}, err
		}

		accountIds = append(accountIds, values[0].(uuid.UUID).String())
		accountsToReturn = append(accountsToReturn, Account{
			ID:      values[0].(uuid.UUID),
			Created: values[1].(time.Time),
			Name:    values[2].(string),
		})
	}

	d.Log.Debug("Accounts retrieved from db", "numAccounts", len(accountsToReturn))

	fieldsSQL := "SELECT \"accountId\", name, value FROM \"accountsFields\" WHERE \"accountId\" = ANY ($1)"
	rows, err = d.DbPool.Query(context.Background(), fieldsSQL, accountIds)
	if err != nil {
		d.Log.Error("Error executing SELECT query for accountFields", "err", err.Error())
		return []Account{}, err
	}

	for rows.Next() {
		var id uuid.UUID
		var name string
		var value []string
		err = rows.Scan(&id, &name, &value)
		if err != nil {
			d.Log.Error("Could not get name or value from database row", "err", err.Error())
			return []Account{}, err
		}

		_, found := accountFields[id]
		if !found {
			accountFields[id] = make(map[string][]string)
		}

		accountFields[id][name] = value
	}

	for idx, account := range accountsToReturn {
		accountsToReturn[idx].Fields = accountFields[account.ID]
	}

	return accountsToReturn, nil
}

func (d Db) AccountUpdateFields(accountID string, fields []AccountCreateInputFields) (Account, error) {
	d.Log.Context = []interface{}{
		"accountID", accountID,
		"fields", fields,
	}

	// Begin database transaction
	conn, err := d.DbPool.Acquire(context.Background())
	if err != nil {
		d.Log.Error("Could not acquire database connection", "err", err.Error())
		return Account{}, err
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		d.Log.Error("Could not begin database transaction", "err", err.Error())
		return Account{}, err
	}

	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), "DELETE FROM \"accountsFields\" WHERE \"accountId\" = $1;", accountID)
	if err != nil {
		d.Log.Error("Could not delete previous fields", "err", err.Error())
		return Account{}, err
	}

	accountFieldsSQL := "INSERT INTO \"accountsFields\" (id, \"accountId\", name, value) VALUES($1,$2,$3,$4);"
	for _, field := range fields {
		newFieldID, err := uuid.NewRandom()
		if err != nil {
			d.Log.Error("Could not create new Uuid", "err", err.Error())
			return Account{}, err
		}

		_, err = tx.Exec(context.Background(), accountFieldsSQL, newFieldID, accountID, field.Name, field.Values)
		if err != nil {
			d.Log.Error("Database error when trying to add account field", "err", err.Error(), "fieldName", field.Name, "fieldvalues", field.Values)
		}

		d.Log.Debug("Added account field", "accountID", accountID, "fieldName", field.Name, "fieldValues", field.Values)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		d.Log.Error("Database error when tying to commit", "err", err.Error())
		return Account{}, err
	}

	return d.AccountGet(accountID, "", "")
}

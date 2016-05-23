package model

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/twinj/uuid"
)

type Account struct {
	ID        string             // uuid
	CreatedAt time.Time
}

func NewAccount() *Account {
	var a Account
	uuid := uuid.NewV4()
	a.ID = uuid.String()
	a.CreatedAt = time.Now()
	return &a
}

func GetAccount(db *bolt.DB, uuid string) (*Account, error) {
	var account Account

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Accounts"))
		v := b.Get([]byte(uuid))

		err := deserialize(&v, &account)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to deserialize object: %s", err)
	}

	return &account, nil
}

// ListAccounts returns all accounts in a map with the uuid as key
func ListAccounts(db *bolt.DB) (*map[string]Account, error) {
	accounts := make(map[string]Account)

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Accounts"))

		return b.ForEach(func(k, v []byte) error {
			if v == nil {
				// v == nil means it's a nested bucket so ignore it
				return nil
			}

			var a Account
			err := deserialize(&v, &a)
			if err != nil {
				glog.Errorf("Failed to deserialize account: %s", err)
				return fmt.Errorf("Failed to deserialize account: %s", err)
			}
			accounts[a.ID] = a

			return nil
		})
	})
	if err != nil {
		glog.Errorf("Failed to get accounts: %s", err)
		return nil, fmt.Errorf("Failed to get accounts: %s", err)
	}

	return &accounts, nil
}

// Save saves the account
func (account *Account) Save(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Accounts"))

		glog.Infof("Saving account %s", account.ID)
		err := BoltSaveObject(b, account.ID, account)
		if err != nil {
			return fmt.Errorf("Failed to save account: %s", err)
		}

		return nil
	})
}

// Alerts returns the alerts for the account
func (account *Account) Alerts(db *bolt.DB) (*map[string]Alert, error) {
	return ListAlerts(db, account.ID)
}

// APIKeys returns all api keys for the account
func (account *Account) APIKeys(db *bolt.DB) (*map[string]APIKey, error) {
	return ListAPIKeys(db, account.ID)
}

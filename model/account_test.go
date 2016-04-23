package model

import (
	"testing"

	"github.com/boltdb/bolt"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	assert := assert.New(t)

	a := NewAccount()
	assert.NotNil(a)
	assert.NotNil(a.ID)
	assert.NotNil(a.CreatedAt)
}

func TestListAccountsWithNoAccount(t *testing.T) {
	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		accounts, err := ListAccounts(db)
		assert.NoError(err)
		assert.NotNil(accounts)
		assert.True(len(*accounts) == 0)
	})
}

func TestAccounts(t *testing.T) {
	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		// should be empty
		accounts, err := ListAccounts(db)
		assert.NoError(err)
		assert.NotNil(accounts)
		assert.True(len(*accounts) == 0)

		// add an account
		a := NewAccount()
		err = a.Save(db)
		assert.NoError(err)

		// should be one
		accounts, err = ListAccounts(db)
		assert.NoError(err)
		assert.NotNil(accounts)
		assert.True(len(*accounts) == 1)
		assert.Equal(*a, (*accounts)[a.ID])

		// add another account
		b := NewAccount()
		err = b.Save(db)
		assert.NoError(err)

		// should be two
		accounts, err = ListAccounts(db)
		assert.NoError(err)
		assert.NotNil(accounts)
		assert.True(len(*accounts) == 2)
		assert.Equal(*a, (*accounts)[a.ID])
		assert.Equal(*b, (*accounts)[b.ID])
	})
}

package model

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestNewAPIKey(t *testing.T) {
	assert := assert.New(t)

	a := NewAPIKey()
	assert.NotNil(a)
	assert.NotNil(a.ID)
	assert.NotNil(a.CreatedAt)
}

func TestAPIKey(t *testing.T) {
	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		// should be empty
		apiKeys, err := ListAPIKeys(db, "foo")
		assert.NoError(err)
		assert.NotNil(apiKeys)
		assert.True(len(*apiKeys) == 0)

		// add a APIKey
		a1 := NewAPIKey()
		err = a1.Save(db, "foo")
		assert.NoError(err)

		// should be one APIKey
		apiKeys, err = ListAPIKeys(db, "foo")
		assert.NoError(err)
		assert.NotNil(apiKeys)
		assert.True(len(*apiKeys) == 1)
		assert.Equal(*a1, (*apiKeys)[a1.ID])

		// add another APIKey
		a2 := NewAPIKey()
		err = a2.Save(db, "foo")
		assert.NoError(err)

		// should be two APIKeys
		apiKeys, err = ListAPIKeys(db, "foo")
		assert.NoError(err)
		assert.NotNil(apiKeys)
		assert.True(len(*apiKeys) == 2)
		assert.Equal(*a1, (*apiKeys)[a1.ID])
		assert.Equal(*a2, (*apiKeys)[a2.ID])

	})
}

package model

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	assert := assert.New(t)

	to := NewToken()
	assert.NotNil(to)
	assert.NotNil(to.ID)
	assert.NotNil(to.CreatedAt)
}

func TestToken(t *testing.T) {
	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		// should be empty
		tokens, err := ListTokens(db, "foo")
		assert.NoError(err)
		assert.NotNil(tokens)
		assert.True(len(*tokens) == 0)

		// add a Token
		t1 := NewToken()

		tt := make(map[string]Token)
		tt[t1.ID] = *t1

		err = SaveTokens(db, "foo", &tt)
		assert.NoError(err)

		// should be one token
		tokens, err = ListTokens(db, "foo")
		assert.NoError(err)
		assert.NotNil(tokens)
		assert.True(len(*tokens) == 1)
		assert.Equal(*t1, (*tokens)[t1.ID])

		//arr another token
		t2 := NewToken()
		tt[t2.ID] = *t2

		err = SaveTokens(db, "foo", &tt)
		assert.NoError(err)

		// should be two tokens
		tokens, err = ListTokens(db, "foo")
		assert.NoError(err)
		assert.NotNil(tokens)
		assert.True(len(*tokens) == 2)
		assert.Equal(*t1, (*tokens)[t1.ID])
		assert.Equal(*t2, (*tokens)[t2.ID])

	})
}

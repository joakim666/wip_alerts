package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerialize(t *testing.T) {
	assert := assert.New(t)

	a := Account{ID: "test"}

	bytes, err := serialize(&a)
	assert.NoError(err)
	assert.True(len(bytes) > 0)
}

func TestDeserialize(t *testing.T) {
	assert := assert.New(t)

	a := Account{ID: "test"}

	b, err := serialize(&a)
	assert.NoError(err)
	assert.True(len(b) > 0)

	var a2 Account

	err = deserialize(&b, &a2)
	assert.NoError(err)
	assert.Equal(a, a2)
}

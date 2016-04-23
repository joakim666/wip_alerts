package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twinj/uuid"
)

func Test(t *testing.T) {
	assert := assert.New(t)
	uuid := uuid.NewV4()
	assert.NotNil(uuid)
}

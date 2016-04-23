package model

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestNewRenewal(t *testing.T) {
	assert := assert.New(t)

	r := NewRenewal()
	assert.NotNil(r)
	assert.NotNil(r.ID)
	assert.NotNil(r.CreatedAt)
}

func TestRenewal(t *testing.T) {
	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		// should be empty
		renewals, err := ListRenewals(db, "foo")
		assert.NoError(err)
		assert.NotNil(renewals)
		assert.True(len(*renewals) == 0)

		// arr a Renewals
		r1 := NewRenewal()

		rr := make(map[string]Renewal)
		rr[r1.ID] = *r1

		err = SaveRenewals(db, "foo", &rr)
		assert.NoError(err)

		// should be one device
		renewals, err = ListRenewals(db, "foo")
		assert.NoError(err)
		assert.NotNil(renewals)
		assert.True(len(*renewals) == 1)
		assert.Equal(*r1, (*renewals)[r1.ID])

		//arr another device
		r2 := NewRenewal()
		rr[r2.ID] = *r2

		err = SaveRenewals(db, "foo", &rr)
		assert.NoError(err)

		// should be two device
		renewals, err = ListRenewals(db, "foo")
		assert.NoError(err)
		assert.NotNil(renewals)
		assert.True(len(*renewals) == 2)
		assert.Equal(*r1, (*renewals)[r1.ID])
		assert.Equal(*r2, (*renewals)[r2.ID])

	})
}

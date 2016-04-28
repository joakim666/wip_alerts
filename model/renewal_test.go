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

		// add a Renewals
		r1 := NewRenewal()
		err = r1.Save(db, "foo")
		assert.NoError(err)

		// should be one renewal
		renewals, err = ListRenewals(db, "foo")
		assert.NoError(err)
		assert.NotNil(renewals)
		assert.True(len(*renewals) == 1)
		assert.Equal(*r1, (*renewals)[r1.ID])

		// add another renewal
		r2 := NewRenewal()
		err = r2.Save(db, "foo")
		assert.NoError(err)

		// should be two renewals
		renewals, err = ListRenewals(db, "foo")
		assert.NoError(err)
		assert.NotNil(renewals)
		assert.True(len(*renewals) == 2)
		assert.Equal(*r1, (*renewals)[r1.ID])
		assert.Equal(*r2, (*renewals)[r2.ID])
	})
}

func TestGetRenewal(t *testing.T) {
	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		// account 'foo' should not have any renewals
		renewals, err := ListRenewals(db, "foo")
		assert.NoError(err)
		assert.NotNil(renewals)
		assert.True(len(*renewals) == 0)

		// account 'bar' should not have any renewals
		renewals, err = ListRenewals(db, "bar")
		assert.NoError(err)
		assert.NotNil(renewals)
		assert.True(len(*renewals) == 0)

		// add a Renewals
		r1 := NewRenewal()
		err = r1.Save(db, "foo")
		assert.NoError(err)

		// account 'foo' should have one renewal
		renewals, err = ListRenewals(db, "foo")
		assert.NoError(err)
		assert.NotNil(renewals)
		assert.True(len(*renewals) == 1)
		assert.Equal(*r1, (*renewals)[r1.ID])

		renewal1, accountID, err := GetRenewal(db, r1.ID)
		assert.NoError(err)
		assert.Equal(r1, renewal1)
		assert.Equal("foo", *accountID)

		// add another renewal
		r2 := NewRenewal()
		err = r2.Save(db, "bar")
		assert.NoError(err)

		// account 'foo' should have one renewal
		renewals, err = ListRenewals(db, "foo")
		assert.NoError(err)
		assert.NotNil(renewals)
		assert.True(len(*renewals) == 1)
		assert.Equal(*r1, (*renewals)[r1.ID])

		// account 'bar' should have one renewal
		renewals, err = ListRenewals(db, "bar")
		assert.NoError(err)
		assert.NotNil(renewals)
		assert.True(len(*renewals) == 1)
		assert.Equal(*r2, (*renewals)[r2.ID])

		renewal2, accountID, err := GetRenewal(db, r2.ID)
		assert.NoError(err)
		assert.Equal(r2, renewal2)
		assert.Equal("bar", *accountID)

	})
}

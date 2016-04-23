package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func RunInTestDb(t *testing.T, f func(t *testing.T, db *bolt.DB)) {
	assert := assert.New(t)

	db, err := newTestDB()
	assert.NoError(err)

	f(t, db)

	closeTestDB(db)
}

func newTestDB() (*bolt.DB, error) {
	// set up a temp path
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}

	path := f.Name()
	f.Close()
	os.Remove(path)

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		buckets := []string{"Accounts", "Devices", "Renewals", "APIKeys", "Heartbeats"}
		for _, b := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(b))
			if err != nil {
				return fmt.Errorf("Failed to create '%s' bucket: %s", b, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func closeTestDB(db *bolt.DB) {
	defer os.Remove(db.Path())
	db.Close()
}

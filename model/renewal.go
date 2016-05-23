package model

import (
	"reflect"
	"time"

	"github.com/boltdb/bolt"
	"github.com/twinj/uuid"
)

type Renewal struct {
	ID             string     // uuid
	RefreshTokenID string     // uuid of refresh token
	UsedAt         *time.Time // the time at which this renewal was used to create a new access token
	CreatedAt      time.Time  // the time at which this renewal was created
}

func (r Renewal) PersistanceID() string {
	return r.ID
}

func (r Renewal) Save(db *bolt.DB, accountUUID string) error {
	return BoltSaveAccountObjects(db, ParentID(accountUUID), "Renewals", BoltSingle(&r))
}

func NewRenewal() *Renewal {
	var r Renewal
	uuid := uuid.NewV4()
	r.ID = uuid.String()
	r.CreatedAt = time.Now()
	return &r
}

func SaveRenewals(db *bolt.DB, accountUUID string, renewals *map[string]Renewal) error {
	return BoltSaveAccountObjects(db, ParentID(accountUUID), "Renewals", BoltMap(renewals))
}

func ListRenewals(db *bolt.DB, accountUUID string) (*map[string]Renewal, error) {
	m, err := BoltGetAccountObjects(db, ParentID(accountUUID), "Renewals", reflect.TypeOf(Renewal{}))
	if err != nil {
		return nil, err
	}

	// convert to map containing Renewal
	m2 := make(map[string]Renewal)
	for _, v := range *m {
		d := v.(*Renewal)
		m2[v.PersistanceID()] = *d
	}

	return &m2, nil
}

// GetRenewal returns the given renewal and accountID if a match is found, nil otherwise
func GetRenewal(db *bolt.DB, renewalID string) (*Renewal, *string, error) {
	o, accountID, err := BoltGetObject(db, "Renewals", renewalID, reflect.TypeOf(Renewal{}))
	if err != nil {
		return nil, nil, err
	}

	if o == nil {
		return nil, nil, nil
	}

	var renewal *Renewal
	renewal = (*o).(*Renewal)
	s := string(*accountID)

	return renewal, &s, nil
}

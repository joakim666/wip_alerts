package model

import (
	"reflect"
	"time"

	"github.com/boltdb/bolt"
	"github.com/twinj/uuid"
)

// APIKeyStatus shows the status of the API Key
type APIKeyStatus string

const (
	// APIKeyActive indicates that this API key is active
	APIKeyActive APIKeyStatus = "active"
	// APIKeyInactive indicates that this API key is inactive
	APIKeyInactive = "inactive"
)

// APIKey contains information about a created API key
type APIKey struct {
	ID          string // uuid
	Description string // the user's description of the key
	Status      APIKeyStatus
	CreatedAt   time.Time
}

// PersistanceID is used by the persistance layer
func (a APIKey) PersistanceID() string {
	return a.ID
}

func (a APIKey) Save(db *bolt.DB, accountUUID string) error {
	return BoltSaveAccountObjects(db, ParentID(accountUUID), "APIKeys", BoltSingle(&a))
}

// NewAPIKey creates a new API key
func NewAPIKey() *APIKey {
	var a APIKey
	uuid := uuid.NewV4()
	a.ID = uuid.String()
	a.Status = APIKeyActive
	a.CreatedAt = time.Now()
	return &a
}

// GetAPIKey returns the API Key with the given id
func GetAPIKey(db *bolt.DB, apiKeyID string) (*APIKey, *string, error) {
	o, parentID, err := BoltGetObject(db, "APIKeys", apiKeyID, reflect.TypeOf(APIKey{}))
	if err != nil {
		return nil, nil, err
	}

	if o == nil {
		return nil, nil, nil
	}

	var apiKey *APIKey
	apiKey = (*o).(*APIKey)
	s := string(*parentID)

	return apiKey, &s, nil
}

// ListAPIKeys returns all API keys for the given account
func ListAPIKeys(db *bolt.DB, accountUUID string) (*map[string]APIKey, error) {
	m, err := BoltGetAccountObjects(db, ParentID(accountUUID), "APIKeys", reflect.TypeOf(APIKey{}))
	if err != nil {
		return nil, err
	}

	// convert to map containing APIKey
	m2 := make(map[string]APIKey)
	for _, v := range *m {
		d := v.(*APIKey)
		m2[v.PersistanceID()] = *d
	}

	return &m2, nil
}

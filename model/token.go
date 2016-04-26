package model

import (
	"reflect"
	"time"

	"github.com/boltdb/bolt"
	"github.com/twinj/uuid"
)

type Scope struct {
	Roles        []string // maps to roles
	Capabilities []string // maps to capabilities
}

type Token struct {
	ID        string    // uuid, maps to jti
	IssueTime time.Time // maps to iat
	Type      string    // maps to type
	Scope     Scope
	RawString string // the raw base64 encoded token data string
	CreatedAt time.Time
}

func (r Token) PersistanceID() string {
	return r.ID
}

func NewToken() *Token {
	var t Token
	uuid := uuid.NewV4()
	t.ID = uuid.String()
	t.CreatedAt = time.Now()
	return &t
}

func SaveTokens(db *bolt.DB, accountUUID string, tokens *map[string]Token) error {
	return BoltSaveAccountObjects(db, accountUUID, "Tokens", BoltMap(tokens))
}

func ListTokens(db *bolt.DB, accountUUID string) (*map[string]Token, error) {
	m, err := BoltGetAccountObjects(db, accountUUID, "Tokens", reflect.TypeOf(Token{}))
	if err != nil {
		return nil, err
	}

	// convert to map containing Token
	m2 := make(map[string]Token)
	for _, v := range *m {
		d := v.(*Token)
		m2[v.PersistanceID()] = *d
	}

	return &m2, nil
}

// ListAllTokens returns a map of account id to array of tokens
func ListAllTokens(db *bolt.DB) (*map[string][]Token, error) {
	m, err := BoltGetObjects(db, "Tokens", reflect.TypeOf(Token{}))
	if err != nil {
		return nil, err
	}

	// convert to map accountId => []Token
	m2 := make(map[string][]Token)
	for k, v := range *m {
		var tokens []Token
		for _, v2 := range v {
			tokens = append(tokens, *v2.(*Token))
		}
		m2[k] = tokens
	}

	return &m2, nil
}
package auth

import (
	"encoding/json"

	"gopkg.in/square/go-jose.v1"
)

type Scope struct {
	Roles        []string `json:"roles"`
	Capabilities []string `json:"capabilities"`
}

type Token struct {
	IssueTime int64  `json:"iat"`
	ID        string `json:"jti"`
	AccountID string `json:"sub"`
	Type      string `json:"type"`
	Scope     Scope  `json:"scope"`
}

// Encrypt the access token using 128 bit AES, with a shared key given in encryptionKey
func EncryptAccessToken(token *Token, encryptionKey interface{}) (string, error) {
	encrypter, err := jose.NewEncrypter(jose.A128KW, jose.A128CBC_HS256, encryptionKey)
	if err != nil {
		return "", err
	}

	return encryptToken(token, encrypter)
}

// Decrypt the access token in serialized form using 128 bit AES, with a shared key given in encryptionKey
func DecryptAccessToken(token string, encryptionKey interface{}) (*Token, error) {
	object, err := jose.ParseEncrypted(token)
	if err != nil {
		return nil, err
	}

	decrypted, err := object.Decrypt(encryptionKey)
	if err != nil {
		return nil, err
	}

	var t Token
	err = json.Unmarshal(decrypted, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Encrypt the refresh token using 2048 bit rsa key. publicKey should be &privateKey.PublicKey
func EncryptRefreshToken(token *Token, publicKey interface{}) (string, error) {
	encrypter, err := jose.NewEncrypter(jose.RSA_OAEP, jose.A128CBC_HS256, publicKey)
	if err != nil {
		return "", err
	}

	return encryptToken(token, encrypter)
}

// Decrypt the refresh token in serialized form using 2048 bit rsa key. privateKey should be &privateKey
func DecryptRefreshToken(token string, privateKey interface{}) (*Token, error) {
	object, err := jose.ParseEncrypted(token)
	if err != nil {
		return nil, err
	}

	decrypted, err := object.Decrypt(privateKey)
	if err != nil {
		return nil, err
	}

	var t Token
	err = json.Unmarshal(decrypted, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Valid validates the access token based on the issue time and expire time (TBD)
func (token *Token) Valid() bool {
	if token.Type != "access_token" {
		return false
	}

	// TODO validate IssueTime

	return true
}

// HasRole checks if the access token contains the given role
func (token *Token) HasRole(role string) bool {
	if len(token.Scope.Roles) == 0 {
		return false
	}

	for _, r := range token.Scope.Roles {
		if r == role {
			return true
		}
	}

	return false
}

// HasCapability checks if the access token contains the given capability
func (token *Token) HasCapability(capability string) bool {
	if len(token.Scope.Capabilities) == 0 {
		return false
	}

	for _, c := range token.Scope.Capabilities {
		if c == capability {
			return true
		}
	}

	return false
}

func encryptToken(token *Token, encrypter jose.Encrypter) (string, error) {
	b, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	object, err := encrypter.Encrypt(b)
	if err != nil {
		return "", err
	}

	serialized, err := object.CompactSerialize()
	if err != nil {
		return "", err
	}

	return serialized, nil
}

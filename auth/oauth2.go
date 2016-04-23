package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// ValidateAccessToken extracts an access token from the headers, checks that it's valid and then passes it on to the check-function
func ValidateAccessToken(check func(token *Token, ctx *gin.Context) bool, encryptionKey interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var serializedToken string
		var err error

		if serializedToken, err = extractToken(c.Request); err != nil {
			glog.Errorf("Can not extract token, caused by: %s", err)
			c.AbortWithError(http.StatusUnauthorized, errors.New("No token supplied"))
			return
		}

		token, err := DecryptAccessToken(serializedToken, encryptionKey)
		if err != nil {
			glog.Errorf("Can not deserialize token, caused by: %s", err)
			c.AbortWithError(http.StatusUnauthorized, errors.New("Token error"))
			return
		}

		if !token.Valid() {
			glog.Error("Token not valid")
			c.AbortWithError(http.StatusUnauthorized, errors.New("Invalid token"))
			return
		}

		if !check(token, c) {
			glog.Errorf("Authorization check failed")
			c.AbortWithError(http.StatusUnauthorized, errors.New("Authorization check failed"))
			return
		}

		glog.Infof("Granting access to %s with roles: %s", token.UserUUID, strings.Join(token.Scope.Roles, ","))

		// access granted
		c.Writer.Header().Set("Bearer", serializedToken)
	}
}

func extractToken(r *http.Request) (string, error) {
	hdr := r.Header.Get("Authorization")
	if hdr == "" {
		return "", errors.New("No authorization header")
	}
	th := strings.Split(hdr, " ")
	if len(th) != 2 {
		return "", errors.New("Incomplete authorization header")
	}

	return th[1], nil
}

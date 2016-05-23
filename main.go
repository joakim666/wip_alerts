package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/joakim666/wip_alerts/auth"
	"errors"
	"net/http"
	"github.com/joakim666/wip_alerts/model"
)

func main() {
	// flag parsing (and setting through code) for glog
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	// Open the my.db data file in the your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Creating buckets")
	glog.Infof("Creating buckets")
	err = db.Update(func(tx *bolt.Tx) error {
		// create all buckets
		buckets := []string{"Accounts", "Devices", "Renewals", "APIKeys", "Heartbeats", "Tokens", "Alerts"}
		for _, b := range buckets {
			glog.Infof("Creating %s bucket", b)
			_, err := tx.CreateBucketIfNotExists([]byte(b))
			if err != nil {
				return fmt.Errorf("Failed to create %s bucket: %s", b, err)
			}
		}

		return nil
	})
	if err != nil {
		glog.Errorf("Bolt failed: %s", err)
	}

	r := setupRoutes(db)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func setupRoutes(db *bolt.DB) *gin.Engine {
	r := gin.Default()

	var sharedKey = []byte("shared key123456") // used for access tokens
	var privateKey = []byte("fooo")            // used for refresh tokens
	var publicKey = []byte("fooo")             // used for refresh tokens

	/* Public routes are as they are named public. No form of authentication or authorization is needed. */
	// Begin: PUBLIC routes
	public := r.Group("/api/v1")
	public.POST("/accounts", PostAccounts(db))
	public.POST("/renewals", PostRenewals(db, privateKey))
	public.POST("/tokens", PostTokens(db, publicKey, sharedKey))
	// End: PUBLIC routes

	/* Api key routes require an api-key, either through a header or as a query-parameter. */
	// Begin: APIKEY routes
	apiKey := r.Group("/api/v1")
	apiKey.Use(validateApiKey(db))
	apiKey.POST("/alerts", CreateAlertRoute(db))
	apiKey.POST("/heartbeats", helloWorld)
	// END: APIKEY routes

	/* Access token routes require an access token set as a header. */
	// Begin: ACCESSTOKEN routes
	private := r.Group("/api/v1")
	private.Use(auth.ValidateAccessToken(hasRole("role1"), sharedKey))
	private.GET("/api-keys", ListAPIKeyRoute(db))
	private.POST("/api-keys", CreateAPIKeyRoute(db))
	private.GET("/ping", PingRoute())
	private.GET("/alerts", helloWorld)
	private.POST("/alerts/:id", helloWorld)
	private.GET("/heartbeats", helloWorld)
	// End: ACCESSTOKEN routes

	/* Admin capability routes requires a token with admin capabilty set */
	// Begin: Admin capability routes
	admin := r.Group("/api/v1")
	admin.Use(auth.ValidateAccessToken(hasCapability("admin"), sharedKey))
	admin.GET("/accounts", ListAccounts(db))
	admin.GET("/renwewals", ListRenewals(db))
	admin.GET("/tokens", ListTokens(db))
	// End: Admin capability routes

	return r
}

func helloWorld(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello world",
	})
}

func hasRole(role string) func(*auth.Token, *gin.Context) bool {
	return func(token *auth.Token, ctx *gin.Context) bool {
		if token.HasRole(role) {
			ctx.Set("accountID", token.AccountID)
			//ctx.Set("account", GetACcount(...)) // TODO fix
			return true
		}

		return false
	}
}

func hasCapability(capability string) func(*auth.Token, *gin.Context) bool {
	return func(token *auth.Token, ctx *gin.Context) bool {
		if token.HasCapability(capability) {
			ctx.Set("accountID", token.AccountID)
			return true
		}

		return false
	}
}

func validateApiKey(db *bolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKeyID, err := extractApiKey(c)
		if err != nil {
			glog.Errorf("Can not find api key: %s", err)
			c.AbortWithError(http.StatusUnauthorized, errors.New("API Key missing"))
			return
		}

		apiKey, accountID, err := model.GetAPIKey(db, apiKeyID)
		if err != nil {
			glog.Errorf("Can not find api key: %s", err)
			c.AbortWithError(http.StatusUnauthorized, errors.New("API Key missing"))
			return
		}

		if model.APIKeyActive != apiKey.Status {
			glog.Errorf("Api key with id %s is not valid", apiKeyID)
			c.AbortWithError(http.StatusUnauthorized, errors.New("API Key not valid"))
			return
		}

		glog.Infof("Granting api level access to api key with id %s", apiKey.ID)
		c.Set("apiKeyID", apiKey.ID)
		c.Set("accountID", accountID)
	}
}

func extractApiKey(c *gin.Context) (string, error) {
	hdr := c.Request.Header.Get("APIKey")
	if hdr != "" {
		return hdr, nil
	}

	q := c.Query("apiKey")
	if q != "" {
		return q, nil
	}

	return "", errors.New("No APIKey found as header or query parameter")
}

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/joakim666/wip_alerts/auth"
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
		buckets := []string{"Accounts", "Devices", "Renewals", "APIKeys", "Heartbeats"}
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

	var sharedKey = []byte("shared key123456")

	public := r.Group("/api/v1")
	public.POST("/accounts", postAccounts(db))
	public.GET("/accounts", listAccounts(db)) // TODO require admin capability
	public.POST("/renewals", helloWorld)
	public.POST("/tokens", helloWorld)

	private := r.Group("/api/v1")
	private.Use(auth.ValidateAccessToken(hasRole("role1"), sharedKey))
	private.GET("/api-keys", helloWorld)
	private.POST("/api-keys", helloWorld)
	private.GET("/ping", helloWorld)
	private.GET("/alerts", helloWorld)
	private.POST("/alerts", helloWorld)
	private.POST("/alerts/:id", helloWorld)
	private.GET("/heartbeats", helloWorld)
	private.POST("/heartbeats", helloWorld)

	/*	private.GET("/ping", func(c *gin.Context) {

			var res []byte

			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("MyBucket"))
				v := b.Get([]byte("answer"))

				res = make([]byte, len(v))
				copy(res, v)

				return nil
			})

			c.JSON(200, gin.H{
				"message": "pong",
				"answer":  string(res),
			})
		})

		tst := r.Group("api")
		private.Use(auth.ValidateAccessToken(hasRole("role2"), sharedKey))
		tst.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "api hello",
			})
		})*/

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
			return true
		}

		return false
	}
}

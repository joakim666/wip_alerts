package main

import (
	"testing"
	"flag"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"net/http"
	"strings"
)

func TestCreateAlertRouteWithMissingAccountID(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.POST("/alerts", CreateAlertRoute(db))

		req, _ := http.NewRequest("POST", "/alerts", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(401, res.Code)
	})
}

func TestCreateAlertRouteWithMissingAPIKeyID(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Set("accountID", "55")
		})

		router.POST("/alerts", CreateAlertRoute(db))

		req, _ := http.NewRequest("POST", "/alerts", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(401, res.Code)
	})
}

func TestCreateAlertRouteWithMissingData(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Set("accountID", "55")
			c.Set("apiKeyID", "55")
		})

		router.POST("/alerts", CreateAlertRoute(db))

		var bodies []string

		// 1. missing all
		bodies = append(bodies, "")

		// 2.
		bodies = append(bodies, `
			{
				"title": "title1"
			}
		`)

		// 3.
		bodies = append(bodies, `
			{
				"title": "title1",
				"short_description": "short_description1"
			}
		`)

		// 4.
		bodies = append(bodies, `
			{
				"title": "title1",
				"short_description": "short_description1",
				"long_description": "long_description"
			}
		`)

		// 5.
		bodies = append(bodies, `
			{
				"title": "title1",
				"short_description": "short_description1",
				"long_description": "long_description",
				"priority": "high"
			}
		`)

		for _, v := range bodies {
			req, _ := http.NewRequest("POST", "/alerts", strings.NewReader(v))
			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)
			assert.Equal(400, res.Code)
		}
	})
}

func TestCreateAlertRouteWithValidData(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Set("accountID", "55")
			c.Set("apiKeyID", "55")
		})

		router.POST("/alerts", CreateAlertRoute(db))

		body := `
			{
				"title": "title1",
				"short_description": "short_description1",
				"long_description": "long_description",
				"priority": "high",
				"triggered_at": "2012-04-23T18:25:43.511Z"
			}
		`

		req, _ := http.NewRequest("POST", "/alerts", strings.NewReader(body))
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(201, res.Code)
	})
}

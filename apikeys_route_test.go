package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"flag"
	"io/ioutil"
	"encoding/json"
	"github.com/joakim666/wip_alerts/model"
)

func TestCreateAPIKeyWithInvalidData(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Set("accountID", "55")
		})

		router.POST("/api-keys", CreateAPIKeyRoute(db))

		var bodies []string

		// 1. missing all
		bodies = append(bodies, "")

		// 2. missing description
		bodies = append(bodies, `
			{}
		`)

		for _, v := range bodies {
			req, _ := http.NewRequest("POST", "/api-keys", strings.NewReader(v))
			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)
			assert.Equal(400, res.Code)
		}
	})
}

func TestCreateAPIKey(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Set("accountID", "55")
		})

		router.POST("/api-keys", CreateAPIKeyRoute(db))

		var body = `
			{
				"description": "new shiny api key"
			}
		`

		req, _ := http.NewRequest("POST", "/api-keys", strings.NewReader(body))
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(201, res.Code)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)

		var resJson interface{}
		err = json.Unmarshal(resBody, &resJson)
		assert.NoError(err)

		resMap := resJson.(map[string]interface{})

		assert.Equal(4, len(resMap))
		assert.NotEmpty(resMap["id"])
		assert.Equal("new shiny api key", resMap["description"])
		assert.NotEmpty(resMap["issued_at"])
		assert.Equal("active", resMap["status"])
	})
}

func TestListAPIKeyWithNoAPIKeysPresent(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Set("accountID", "55")
		})

		router.GET("/api-keys", ListAPIKeyRoute(db))

		req, _ := http.NewRequest("GET", "/api-keys", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(200, res.Code)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)

		var resJson interface{}
		err = json.Unmarshal(resBody, &resJson)
		assert.NoError(err)

		resMap := resJson.(map[string]interface{})
		assert.Equal(0, len(resMap)) // should be empty map
	})
}

func TestListAPIKey(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		apiKey1 := model.NewAPIKey()
		apiKey1.Description = "my description"
		apiKey1.Save(db, "55")

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Set("accountID", "55")
		})

		router.GET("/api-keys", ListAPIKeyRoute(db))

		req, _ := http.NewRequest("GET", "/api-keys", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(200, res.Code)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)

		var resJson interface{}
		err = json.Unmarshal(resBody, &resJson)
		assert.NoError(err)

		resMap := resJson.(map[string]interface{})
		assert.Equal(1, len(resMap)) // should not be empty return

		r1 := resMap[apiKey1.ID].(map[string]interface{})

		assert.Equal(apiKey1.ID, r1["id"])
		assert.Equal(apiKey1.Description, r1["description"])
		assert.Equal(string(apiKey1.Status), r1["status"])
		assert.NotEmpty(r1["issued_at"])

		apiKey2 := model.NewAPIKey()
		apiKey2.Description = "my description"
		apiKey2.Save(db, "55")
	})
}

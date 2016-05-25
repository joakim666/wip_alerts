package main

import (
	"flag"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"io/ioutil"
	"github.com/joakim666/wip_alerts/model"
	"time"
	"testing"
	"net/http"
	"encoding/json"
)

func TestCreateHeartbeatRouteWithMissingAccountID(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.POST("/heartbeats", CreateHeartbeatRoute(db))

		req, _ := http.NewRequest("POST", "/heartbeats", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(http.StatusUnauthorized, res.Code)
	})
}

func TestCreateHeartbeatRouteWithMissingAPIKeyID(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Set("accountID", "55")
		})

		router.POST("/heartbeats", CreateHeartbeatRoute(db))

		req, _ := http.NewRequest("POST", "/heartbeats", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(http.StatusUnauthorized, res.Code)
	})
}

func TestCreateHeartbeatRouteWithMissingData(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Set("accountID", "55")
			c.Set("apiKeyID", "55")
		})

		router.POST("/heartbeats", CreateHeartbeatRoute(db))

		var bodies []string

		// 1. missing all
		bodies = append(bodies, "")

		// 2.
		bodies = append(bodies, `
			{
				"foo": "title1"
			}
		`)

		for _, v := range bodies {
			req, _ := http.NewRequest("POST", "/heartbeats", strings.NewReader(v))
			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)
			assert.Equal(http.StatusBadRequest, res.Code)
		}
	})
}

func TestCreateHeartbeatRouteWithValidData(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Set("accountID", "55")
			c.Set("apiKeyID", "55")
		})

		router.POST("/heartbeats", CreateHeartbeatRoute(db))

		body := `
			{
				"executed_at": "2012-04-23T18:25:43.511Z"
			}
		`

		req, _ := http.NewRequest("POST", "/heartbeats", strings.NewReader(body))
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)

		var resJson interface{}
		err = json.Unmarshal(resBody, &resJson)
		assert.NoError(err)

		resMap := resJson.(map[string]interface{})
		assert.NotEmpty(resMap["id"])
		assert.NotEmpty(resMap["created_at"])
		assert.Equal("2012-04-23T18:25:43.511Z", resMap["executed_at"])
	})
}

func TestLatestHeartbeatsWithoutAccountID(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.GET("/heartbeats", LatestHeartbeatsRoute(db))

		req, _ := http.NewRequest("GET", "/heartbeats", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(http.StatusUnauthorized, res.Code)
	})
}

func TestLatestHeartbeats(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Set("accountID", "55")
		})

		router.GET("/heartbeats", LatestHeartbeatsRoute(db))

		req, _ := http.NewRequest("GET", "/heartbeats", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)
		resBody, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)

		var resJson interface{}
		err = json.Unmarshal(resBody, &resJson)
		assert.NoError(err)

		resMap := resJson.(map[string]interface{})
		// Assert 0 heartbeats for this account
		assert.Equal(0, len(resMap))

		apiKey1 := model.NewAPIKey()
		apiKey1.Description = "my description"
		err = apiKey1.Save(db, "55")
		assert.NoError(err)

		// add an alert for apiKey1
		h1 := model.NewHeartbeat(apiKey1.ID)
		h1.ExecutedAt = time.Now()
		err = h1.Save(db, "55")
		assert.NoError(err)

		req, _ = http.NewRequest("GET", "/heartbeats", nil)
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)
		resBody, err = ioutil.ReadAll(res.Body)
		assert.NoError(err)

		err = json.Unmarshal(resBody, &resJson)
		assert.NoError(err)

		resMap = resJson.(map[string]interface{})
		// Assert 1 heartbeat for this account
		assert.Equal(1, len(resMap))

		r1 := resMap[h1.ID].(map[string]interface{})
		assert.NotEmpty(r1["executed_at"])

		// add another api key for the same account
		apiKey2 := model.NewAPIKey()
		apiKey2.Description = "my description2"
		err = apiKey2.Save(db, "55")
		assert.NoError(err)

		// add an heartbeat for apiKey2
		h2 := model.NewHeartbeat(apiKey2.ID)
		h2.ExecutedAt = time.Now().Add(-10 * time.Minute)
		err = h2.Save(db, "55")
		assert.NoError(err)

		req, _ = http.NewRequest("GET", "/heartbeats", nil)
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)
		resBody, err = ioutil.ReadAll(res.Body)
		assert.NoError(err)

		err = json.Unmarshal(resBody, &resJson)
		assert.NoError(err)

		resMap = resJson.(map[string]interface{})
		// Assert 2 heartbeats for this account
		assert.Equal(2, len(resMap))

		r1 = resMap[h1.ID].(map[string]interface{})
		assert.NotEmpty(r1["executed_at"])

		r2 := resMap[h2.ID].(map[string]interface{})
		assert.NotEmpty(r2["executed_at"])

		// add another heartbeat for apiKey2
		h3 := model.NewHeartbeat(apiKey2.ID)
		h3.ExecutedAt = time.Now()
		err = h3.Save(db, "55")
		assert.NoError(err)

		req, _ = http.NewRequest("GET", "/heartbeats", nil)
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)
		resBody, err = ioutil.ReadAll(res.Body)
		assert.NoError(err)

		err = json.Unmarshal(resBody, &resJson)
		assert.NoError(err)

		resMap = resJson.(map[string]interface{})

		// Assert still only two heartbeats for this account, but should be h1 and h3 this time
		assert.Equal(2, len(resMap))

		r1 = resMap[h1.ID].(map[string]interface{})
		assert.NotEmpty(r1["executed_at"])

		r2 = resMap[h3.ID].(map[string]interface{})
		assert.NotEmpty(r2["executed_at"])
	})
}

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
	"io/ioutil"
	"encoding/json"
	"github.com/joakim666/wip_alerts/model"
	"time"
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
				"long_description": "long_description1",
				"priority": "high",
				"triggered_at": "2012-04-23T18:25:43.511Z"
			}
		`

		req, _ := http.NewRequest("POST", "/alerts", strings.NewReader(body))
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(201, res.Code)

		resBody, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)

		var resJson interface{}
		err = json.Unmarshal(resBody, &resJson)
		assert.NoError(err)

		resMap := resJson.(map[string]interface{})
		assert.NotEmpty(resMap["id"])
		assert.Equal("title1", resMap["title"])
		assert.Equal("short_description1", resMap["short_description"])
		assert.Equal("long_description1", resMap["long_description"])
		assert.Equal("high", resMap["priority"])
		assert.Equal("new", resMap["status"])
		assert.NotEmpty(resMap["triggered_at"])
		assert.NotEmpty(resMap["created_at"])
	})
}

func TestListAlertsWithoutAccountID(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.GET("/alerts", ListAlertsRoute(db))

		req, _ := http.NewRequest("GET", "/alerts", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(401, res.Code)
	})
}

func TestListAlerts(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Set("accountID", "55")
		})

		router.GET("/alerts", ListAlertsRoute(db))

		req, _ := http.NewRequest("GET", "/alerts", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(200, res.Code)
		resBody, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)

		var resJson interface{}
		err = json.Unmarshal(resBody, &resJson)
		assert.NoError(err)

		resMap := resJson.(map[string]interface{})
		// Assert 0 alerts for this account
		assert.Equal(0, len(resMap))

		apiKey1 := model.NewAPIKey()
		apiKey1.Description = "my description"
		apiKey1.Save(db, "55")

		// add an alert for apiKey1
		a1 := model.NewAlert(apiKey1.ID)
		a1.Title = "title1"
		a1.ShortDescription = "short_description1";
		a1.LongDescription = "long_description1";
		a1.Priority = model.HighPriority
		a1.TriggeredAt = time.Now()
		a1.Save(db, "55")

		req, _ = http.NewRequest("GET", "/alerts", nil)
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(200, res.Code)
		resBody, err = ioutil.ReadAll(res.Body)
		assert.NoError(err)

		err = json.Unmarshal(resBody, &resJson)
		assert.NoError(err)

		resMap = resJson.(map[string]interface{})
		// Assert 1 alert for this account
		assert.Equal(1, len(resMap))

		r1 := resMap[a1.ID].(map[string]interface{})
		assert.Equal("title1", r1["title"])
		assert.Equal("short_description1", r1["short_description"])
		assert.Equal("long_description1", r1["long_description"])
		assert.Equal("high", r1["priority"])
		assert.Equal("new", r1["status"])
		assert.NotEmpty(r1["triggered_at"])
		assert.NotEmpty(r1["created_at"])

		// adda nother api key for the same account
		apiKey2 := model.NewAPIKey()
		apiKey2.Description = "my description2"
		apiKey2.Save(db, "55")

		// add an alert for apiKey2
		a2 := model.NewAlert(apiKey2.ID)
		a2.Title = "title2"
		a2.ShortDescription = "short_description2";
		a2.LongDescription = "long_description2";
		a2.Priority = model.HighPriority
		a2.TriggeredAt = time.Now()
		a2.Save(db, "55")

		req, _ = http.NewRequest("GET", "/alerts", nil)
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(200, res.Code)
		resBody, err = ioutil.ReadAll(res.Body)
		assert.NoError(err)

		err = json.Unmarshal(resBody, &resJson)
		assert.NoError(err)

		resMap = resJson.(map[string]interface{})
		// Assert 2 alerts for this account
		assert.Equal(2, len(resMap))

		r1 = resMap[a1.ID].(map[string]interface{})
		assert.Equal("title1", r1["title"])
		assert.Equal("short_description1", r1["short_description"])
		assert.Equal("long_description1", r1["long_description"])
		assert.Equal("high", r1["priority"])
		assert.Equal("new", r1["status"])
		assert.NotEmpty(r1["triggered_at"])
		assert.NotEmpty(r1["created_at"])

		r2 := resMap[a2.ID].(map[string]interface{})
		assert.Equal("title2", r2["title"])
		assert.Equal("short_description2", r2["short_description"])
		assert.Equal("long_description2", r2["long_description"])
		assert.Equal("high", r2["priority"])
		assert.Equal("new", r2["status"])
		assert.NotEmpty(r2["triggered_at"])
		assert.NotEmpty(r2["created_at"])

		// adda nother api key for the another account
		apiKey3 := model.NewAPIKey()
		apiKey3.Description = "my description2"
		apiKey3.Save(db, "66")

		// add an alert for apiKey3
		a3 := model.NewAlert(apiKey3.ID)
		a3.Title = "title3"
		a3.ShortDescription = "short_description3";
		a3.LongDescription = "long_description3";
		a3.Priority = model.HighPriority
		a3.TriggeredAt = time.Now()
		a3.Save(db, "66")

		req, _ = http.NewRequest("GET", "/alerts", nil)
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(200, res.Code)
		resBody, err = ioutil.ReadAll(res.Body)
		assert.NoError(err)

		err = json.Unmarshal(resBody, &resJson)
		assert.NoError(err)

		resMap = resJson.(map[string]interface{})

		// Assert still only two alerts for this account
		assert.Equal(2, len(resMap))

		// Change status of alert 1
		a1.Status = model.ArchivedStatus
		a1.Save(db, "55")

		req, _ = http.NewRequest("GET", "/alerts", nil)
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(200, res.Code)
		resBody, err = ioutil.ReadAll(res.Body)
		assert.NoError(err)

		err = json.Unmarshal(resBody, &resJson)
		assert.NoError(err)

		resMap = resJson.(map[string]interface{})
		// Assert 1 alert with status new for this account
		assert.Equal(1, len(resMap))

		// should be alert2 that is returned
		r1 = resMap[a2.ID].(map[string]interface{})
		assert.Equal("title2", r1["title"])
		assert.Equal("short_description2", r1["short_description"])
		assert.Equal("long_description2", r1["long_description"])
		assert.Equal("high", r1["priority"])
		assert.Equal("new", r1["status"])
		assert.NotEmpty(r1["triggered_at"])
		assert.NotEmpty(r1["created_at"])
	})
}

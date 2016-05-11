package main

import (
	"testing"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
)

// TestPing verifies that the ping route returns a 204 status code.
func TestPing(t *testing.T) {
	assert := assert.New(t)
	flag.Lookup("logtostderr").Value.Set("true")

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/ping", PingRoute())

	req, _ := http.NewRequest("GET", "/ping", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	assert.Equal(204, res.Code)
}


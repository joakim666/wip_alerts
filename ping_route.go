package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// PingRoute does nothing but returns a 204 No Content answer. The purpose of this route is to check if an
// access token is still valid.
func PingRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		glog.Infof("PingRoute")

		c.Status(204)
	}
}

package env

import (
	"github.com/mercor/payment-service/constants"
	"github.com/gin-gonic/gin"
)

// Middleware adds env in ctx
func Middleware(e interface{}) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Set(constants.Env, e)
		c.Next()
	}
}

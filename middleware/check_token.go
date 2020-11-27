package middleware

import (
	"InformationPush/controllers"
	"github.com/gin-gonic/gin"
)

/**
* token middleware example
 */
func TokenCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		appId := c.Request.Header.Get("appid")
		token := c.Request.Header.Get("token")
		if appId == "" || token == "" {
			c.Abort()
			controllers.Response(c, 403, "Headers Token does not exist!", nil)
			return
		}
		c.Set("token", token)
		c.Next()
	}
}

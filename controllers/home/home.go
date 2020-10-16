package home

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func Index(c *gin.Context) {
	data := gin.H{
		"title":        "首页",
		"httpUrl":      viper.GetString("app.httpUrl"),
		"webSocketUrl": viper.GetString("app.webSocketUrl"),
	}

	c.HTML(http.StatusOK, "index.html", data)
}

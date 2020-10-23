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
		"miToken":      viper.GetString("test.token.mi"),
		"miAppId":      viper.GetString("test.appId.mi"),
		"hwToken":      viper.GetString("test.token.hw"),
		"hwAppId":      viper.GetString("test.appId.hw"),
		"iosToken":     viper.GetString("test.token.ios"),
		"iosAppId":     viper.GetString("test.appId.ios"),
	}

	c.HTML(http.StatusOK, "index.html", data)
}

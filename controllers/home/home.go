package home

import (
	"InformationPush/common"
	"InformationPush/controllers"
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

func GetAppId(c *gin.Context) {
	channel := c.Param("channel")
	appId := viper.GetString("test.appId." + channel)
	token := viper.GetString("test.token." + channel)
	data := make(map[string]interface{})
	data["appId"] = appId
	data["token"] = token
	controllers.Response(c, common.HTTPOK, "发送成功！", data)
}

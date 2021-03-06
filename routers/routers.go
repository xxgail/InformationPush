package routers

import (
	"InformationPush/controllers/device"
	"InformationPush/controllers/home"
	"InformationPush/controllers/push"
	"InformationPush/middleware"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.LoadHTMLGlob("views/**/*")

	deviceRouter := router.Group("/device")
	{
		deviceRouter.POST("/register", device.Register)
		deviceRouter.POST("/login", device.Login)
		deviceRouter.POST("/logout", device.Logout)
	}

	homeRouter := router.Group("/home")
	{
		homeRouter.GET("/index", home.Index)
		homeRouter.GET("/getAppId/:channel", home.GetAppId)
	}

	pushRouter := router.Group("/push").Use(middleware.TokenCheck()) // middleware example
	{
		pushRouter.POST("/message", push.Message)
	}
}

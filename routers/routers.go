package routers

import (
	"InformationPush/controllers/device"
	"InformationPush/controllers/home"
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
	}
}

//func WebsocketInit() {
//	websocket.Register("addGroup", websocket.AddGroupController)
//	websocket.Register("heartbeat", websocket.HeartbeatController)
//}

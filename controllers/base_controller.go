package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"InformationPush/common"
)

type BaseController struct {
	gin.Context
}

func Response(c *gin.Context, code uint32, msg string, data map[string]interface{}) {

	if code == common.OK {
		code = common.HTTPOK
	} else if code == common.Error {
		code = common.HTTPError
	}

	message := common.Response(code, msg, data)

	// 允许跨域
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") // 服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma, appid, channel, token")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
	c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                                                   //  跨域请求是否需要带cookie信息 默认设置为true
	c.Set("content-type", "application/json")

	c.JSON(http.StatusOK, message)

	return
}

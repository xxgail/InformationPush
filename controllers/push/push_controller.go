package push

import (
	"InformationPush/common"
	"InformationPush/controllers"
	"InformationPush/lib/mysqllib"
	"InformationPush/lib/redislib"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/xxgail/PushSDK"
	"reflect"
	"strconv"
	"time"
)

type MessageParam struct {
	GroupId      string `form:"group_id" json:"group_id"`
	Uid          string `form:"uid" json:"uid"`
	Title        string `form:"title" json:"title" binding:"required"`
	Desc         string `form:"desc" json:"desc" binding:"required"`
	Icon         string `form:"icon" json:"icon"`
	Type         int    `form:"type" json:"type"`
	IsShowNotify int    `form:"is_show_notify" json:"is_show_notify"`
}

var ctx = context.Background()

func Message(c *gin.Context) {
	// 定义接口返回data
	data := make(map[string]interface{})
	contentType := c.Request.Header.Get("Content-Type")
	fmt.Println("Content-Type:", contentType)

	var param MessageParam
	var err error

	switch contentType {
	case "application/json":
		err = c.ShouldBindJSON(&param)
	case "application/x-www-form-urlencoded":
		err = c.ShouldBindWith(&param, binding.Form)
	}
	if err != nil {
		fmt.Println(err)
	}

	channel := c.Request.Header.Get("Channel")
	appId := c.Request.Header.Get("AppId")
	uid := param.Uid

	// 消息体
	send := PushSDK.NewSend()
	send.SetChannel(channel)
	send.SetTitle(param.Title).SetContent(param.Desc)

	apns := common.SMD5(strconv.FormatInt(time.Now().Unix(), 10), "")
	apnsId := apns[:8] + "-" + apns[8:12] + "-" + apns[12:16] + "-" + apns[16:20] + "-" + apns[20:]
	send.SetApnsId(apnsId)

	// 单个查询用户
	//var DeviceToken string
	//mysqlClient := mysqllib.GetMysqlConn()
	//query := "SELECT device_token FROM device WHERE channel = '" + channel + "' AND uid = '" + uid + "' AND app_id = '" + appId + "'"
	//fmt.Println(query)
	//err = mysqlClient.QueryRow(query).Scan(&DeviceToken)
	//if err != nil {
	//	fmt.Println("PushMessage 查询数据库单条用户信息 出错：", err)
	//}

	// 多个查询，获取要推送的device_token
	var pushIds []string
	mysqlClient := mysqllib.GetMysqlConn()
	query := "SELECT device_token FROM device WHERE channel = '" + channel + "' AND uid = '" + uid + "' AND app_id = '" + appId + "'"
	rows, err := mysqlClient.Query(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var pushId string
		if err := rows.Scan(&pushId); err != nil {
			fmt.Println(err)
			return
		}
		pushIds = append(pushIds, pushId)
	}
	if len(pushIds) == 0 {
		controllers.Response(c, common.HTTPError, "No pushId", data)
		return
	}
	send.SetPushId(pushIds)

	// 直接传入配置json
	var plat string
	appIdKey := "AppIdKey:" + channel + ":" + common.SMD5("AppId", appId)
	redisClient := redislib.GetClient()
	platRedis, err := redisClient.HGetAll(ctx, appIdKey).Result()
	if err != nil {
		fmt.Println("err", err)
	} else if len(platRedis) == 0 {
		fmt.Println(channel + "platform param first in redis info")
		query1 := "SELECT value FROM platform WHERE app_id = '" + "1" + "'" + "AND channel = '" + channel + "'"
		err = mysqlClient.QueryRow(query1).Scan(&plat)
		if err != nil {
			fmt.Println("Platform 参数 查询数据库单条用户信息 出错：", err)
			controllers.Response(c, common.HTTPError, "no platform parma", data)
			return
		}
		if plat == "" {
			controllers.Response(c, common.HTTPError, "platform parma is Empty", data)
			return
		}
		// 遍历json，存入redis
		var m map[string]string
		_ = json.Unmarshal([]byte(plat), &m)
		for k, v := range m {
			redisClient.HSet(ctx, appIdKey, k, v)
		}
	} else {
		fmt.Println("get platform param from redis", platRedis, reflect.TypeOf(platRedis))
		// map 转 string
		platStr, _ := json.Marshal(platRedis)
		plat = string(platStr)
	}
	switch channel {
	case "hw":
		send.SetHWParam(plat)
		break
	case "ios":
		send.SetIOSParam(plat)
		break
	case "mi":
		send.SetMIParam(plat)
		break
	case "mz":
		send.SetMZParam(plat)
		break
	case "oppo":
		send.SetOPPOParam(plat)
		break
	case "vivo":
		send.SetVIVOParam(plat)
		break
	}
	if send.Err != nil {
		controllers.Response(c, common.HTTPError, send.Err.Error(), data)
		return
	}
	// 发送消息
	fmt.Println(send)
	response, err := send.SendMessage()
	if err != nil {
		controllers.Response(c, common.HTTPError, "err!", data)
	} else {
		if response.Code == PushSDK.SendSuccess {
			controllers.Response(c, common.HTTPOK, "发送成功！", data)
		} else {
			controllers.Response(c, common.HTTPError, response.Reason, data)
		}
	}
}

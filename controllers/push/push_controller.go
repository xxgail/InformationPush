package push

import (
	"InformationPush/common"
	"InformationPush/controllers"
	"InformationPush/lib/mysqllib"
	"InformationPush/lib/redislib"
	"InformationPush/worker"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"reflect"
)

type MessageParam struct {
	GroupId      string `form:"group_id" json:"group_id"`
	Uid          string `form:"uid" json:"uid"`
	Title        string `form:"title" json:"title" binding:"required"`
	Desc         string `form:"desc" json:"desc" binding:"required"`
	Icon         string `form:"icon" json:"icon"`
	Type         int    `form:"type" json:"type"`
	IsShowNotify int    `form:"is_show_notify" json:"is_show_notify"`
	SendTime     string `form:"send_time" json:"send_time"`
}

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

	appId := c.Request.Header.Get("AppId")
	uid := param.Uid
	uidStr := common.ToSqlStr(uid, ",")

	mysqlClient := mysqllib.GetMysqlConn()
	query := "SELECT device_token, channel FROM device WHERE app_id = '" + appId + "' AND uid IN (" + uidStr + ")" + " LIMIT 0, 1000"
	fmt.Println(query)
	rows, err := mysqlClient.Query(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	pushMap := map[string]string{}
	for rows.Next() {
		var pushId string
		var pushChannel string
		if err := rows.Scan(&pushId, &pushChannel); err != nil {
			fmt.Println(err)
			return
		}
		if _, ok := pushMap[pushChannel]; ok {
			pushMap[pushChannel] = pushMap[pushChannel] + "," + pushId
		} else {
			pushMap[pushChannel] = pushId
		}
	}

	if len(pushMap) == 0 {
		controllers.Response(c, common.HTTPError, "No pushId", data)
		return
	}

	for k, v := range pushMap {
		var plat string
		appIdKey := "AppIdKey:" + k + ":" + common.SMD5("AppId", appId)
		redisClient := redislib.GetClient()
		platRedis, err := redisClient.HGetAll(common.Ctx, appIdKey).Result()
		if err != nil {
			fmt.Println("err", err)
		} else if len(platRedis) == 0 {
			fmt.Println(k + "platform param first in redis info")
			query1 := "SELECT value FROM platform WHERE app_id = '" + "1" + "'" + "AND channel = '" + k + "'"
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
				redisClient.HSet(common.Ctx, appIdKey, k, v)
			}
		} else {
			fmt.Println("get platform param from redis", platRedis, reflect.TypeOf(platRedis))
			// map 转 string
			platStr, _ := json.Marshal(platRedis)
			plat = string(platStr)
		}

		messageTask := worker.MessageTaskStruct{
			Title:   param.Title,
			Content: param.Desc,
			PushId:  v,
			Channel: k,
			Plat:    plat,
		}
		messageTaskStr, _ := json.Marshal(messageTask)
		pushJobParam := worker.PushJobParam{
			SendMessageParam: string(messageTaskStr),
			DelayTime:        param.SendTime,
		}
		sendStr, _ := json.Marshal(pushJobParam)
		res, _ := redisClient.RPush(common.Ctx, "SendMessage", string(sendStr)).Result()
		fmt.Println("存入redis", res)
	}

	controllers.Response(c, common.HTTPOK, "发送成功！", data)
}

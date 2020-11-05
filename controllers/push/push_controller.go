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
	"github.com/mitchellh/mapstructure"
	"github.com/xxgail/PushSDK"
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

	message := PushSDK.MessageBody{
		Title: param.Title,
		Desc:  param.Desc,
	}

	//var DeviceToken string
	//mysqlClient := mysqllib.GetMysqlConn()
	//query := "SELECT device_token FROM device WHERE channel = '" + channel + "' AND uid = '" + uid + "' AND app_id = '" + appId + "'"
	//fmt.Println(query)
	//err = mysqlClient.QueryRow(query).Scan(&DeviceToken)
	//if err != nil {
	//	fmt.Println("PushMessage 查询数据库单条用户信息 出错：", err)
	//}

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

	var plat PushSDK.PlatformParam
	appIdKey := "AppIdKey:" + common.SMD5("AppId", appId)
	fmt.Println("appIdKey ", appIdKey)
	redisClient := redislib.GetClient()
	platRedis, err := redisClient.HGetAll(ctx, appIdKey).Result()
	if err != nil {
		fmt.Println("err", err)
	} else if len(platRedis) == 0 {
		fmt.Println("first in redis info")
		query1 := "SELECT hw_appId, hw_clientSecret, iOS_keyId, iOS_teamId, iOS_bundleId, iOS_authTokenPath, mi_appSecret, mi_restrictedPackageName, mz_appId, mz_appSecret, oppo_appKey, oppo_masterSecret FROM platform_param WHERE app_id = '" + "1" + "'"
		fmt.Println(query1)
		err = mysqlClient.QueryRow(query1).Scan(&plat.HWAppId, &plat.HWClientSecret, &plat.IOSKeyId, &plat.IOSTeamId, &plat.IOSBundleId, &plat.IOSAuthTokenPath, &plat.MIAppSecret, &plat.MIRestrictedPackageName, &plat.MZAppId, &plat.MZAppSecret, &plat.OPPOAppKey, &plat.OPPOMasterSecret)
		if err != nil {
			fmt.Println("Platform 参数 查询数据库单条用户信息 出错：", err)
		}
		platStr, err := json.Marshal(plat)
		if err != nil {

		}
		fmt.Println("platStr", string(platStr))
		var platMap map[string]string
		err = json.Unmarshal(platStr, &platMap)
		if err != nil {

		}
		fmt.Println("platMap", platMap)
		for k, v := range platMap {
			redisClient.HSet(ctx, appIdKey, k, v)
		}
	} else {
		fmt.Println("get plat from redis")
		if err = mapstructure.Decode(platRedis, &plat); err != nil {
			fmt.Println(err)
		}
	}
	send := PushSDK.InitSend(message, channel, pushIds, plat)
	fmt.Println(send)
	code, reason := send.SendMessage()

	if code == 1 {
		controllers.Response(c, common.HTTPOK, "发送成功！", data)
	} else {
		data["reason"] = reason
		controllers.Response(c, common.HTTPError, "api错误", data)
	}
}

//消息payload，根据业务自定义
//type Payload struct {
//	PushTitle    string `json:"push_title"`
//	PushBody     string `json:"push_body"`
//	IsShowNotify string `json:"is_show_notify"`
//	Ext          string `json:"ext"`
//}

//func miGroupPush(regIds []string, payload *Payload, appId string) int {
//	appSecret := viper.GetString("mi.appSecret")
//	restrictedPackageName := viper.GetString("mi.restrictedPackageName")
//	payloadStr, _ := json.Marshal(payload)
//	//是否透传
//	passThrough := "1"
//	if payload.IsShowNotify == "1" {
//		passThrough = "0" //通知栏消息
//	}
//
//	var message = XMPushSDK.InitMessage(payload.PushTitle, payload.PushBody, restrictedPackageName, string(payloadStr), passThrough)
//
//	result, err := XMPushSDK.MessageSend(appSecret, message, regIds)
//	fmt.Println(result)
//	if err != nil {
//		fmt.Println("群发推送报错：", err)
//	}
//	if result != nil && result.Code != XMPushSDK.Success {
//		fmt.Println("群发推送失败，失败原因：", result.Description)
//		return 0
//	}
//	return 1
//}
//
//func hwPush(tokens []string, payload *Payload, appId string) int {
//	clientSecret := viper.GetString("hw.clientSecret")
//	//restrictedPackageName := viper.GetString("hw.restrictedPackageName")
//	//payloadStr, _ := json.Marshal(payload)
//	//是否透传
//	passThrough := "1"
//	if payload.IsShowNotify == "1" {
//		passThrough = "0" //通知栏消息
//	}
//	var message = HWPushSDK.InitMessage(payload.PushTitle, payload.PushBody, passThrough, tokens)
//	fmt.Println("message:----", message)
//	result, err := HWPushSDK.MessagesSend(message, appId, clientSecret)
//	fmt.Println(result)
//	if err != nil {
//		fmt.Println("群发推送报错：", err)
//	}
//	if result != nil && result.Code != HWPushSDK.Success {
//		fmt.Println("群发推送失败，失败原因：", result.Msg)
//		return 0
//	}
//	return 1
//}
//
//func iOSPush(regId string, payload *Payload, appId string) int {
//	keyId := viper.GetString("ios.keyId")
//	teamId := viper.GetString("ios.teamId")
//	bundleID := viper.GetString("ios.bundleID")
//	fmt.Println(keyId, teamId)
//	//是否透传
//	passThrough := "1"
//	if payload.IsShowNotify == "1" {
//		passThrough = "0" //通知栏消息
//	}
//	var message = iOSPushSDK.InitMessage(payload.PushTitle, payload.PushBody, passThrough)
//	authToken, err := iOSPushSDK.GetAuthToken("./config/iosP8/AuthKey_R66FMTN5B2.p8", keyId, teamId)
//	if err != nil {
//		log.Panicln(err)
//	}
//	fmt.Println("aaaaaa:", authToken)
//	result, err := iOSPushSDK.MessagesSend(message, regId, authToken, bundleID)
//	fmt.Println(result)
//	if err != nil {
//		fmt.Println("群发推送报错：", err)
//	}
//	if result != nil && result.Status != iOSPushSDK.Success {
//		fmt.Println("群发推送失败，失败原因：", result.Reason)
//		return 0
//	}
//	return 1
//}
//
//func mzPush(regId []string, payload *Payload, appId string) int {
//	appSecret := viper.GetString("mz.appSecret")
//	fmt.Println(appId, appSecret)
//
//	var message = MZPushSDK.InitMessage(payload.PushTitle, payload.PushBody)
//	result, err := MZPushSDK.MessageSend(message, appId, regId, appSecret)
//	fmt.Println(result)
//	if err != nil {
//		fmt.Println("群发推送报错：", err)
//	}
//	if result != nil && result.Code != MZPushSDK.Success {
//		fmt.Println("群发推送失败，失败原因：", result.Message)
//		return 0
//	}
//	return 1
//}
//
//func oppoPush(regId []string, payload *Payload, appId string) int {
//	appKey := viper.GetString("oppo.appKey")
//	masterSecret := viper.GetString("oppo.masterSecret")
//	fmt.Println(appKey, masterSecret)
//
//	var message = OPPOPushSDK.InitMessage(payload.PushTitle, payload.PushBody, regId)
//	authToken, err := OPPOPushSDK.GetAuthToken(appKey, masterSecret)
//	if err != nil {
//		log.Panicln(err)
//	}
//	result, err := OPPOPushSDK.MessageSend(message, authToken)
//	fmt.Println(result)
//	if err != nil {
//		fmt.Println("群发推送报错：", err)
//	}
//	if result != nil && result.Code != OPPOPushSDK.Success {
//		fmt.Println("群发推送失败，失败原因：", result.Message)
//		return 0
//	}
//	return 1
//}

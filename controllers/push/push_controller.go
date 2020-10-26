package push

import (
	"InformationPush/common"
	"InformationPush/controllers"
	"InformationPush/lib/mysqllib"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/viper"
	"github.com/xxgail/HWPushSDK"
	"github.com/xxgail/MZPushSDK"
	"github.com/xxgail/XMPushSDK"
	"github.com/xxgail/iOSPushSDK"
	"log"
	"strconv"
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

func Message(c *gin.Context) {
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

	var DeviceToken string
	mysqlClient := mysqllib.GetMysqlConn()
	query := "SELECT device_token FROM device WHERE channel = '" + channel + "' AND uid = '" + uid + "' AND app_id = '" + appId + "'"
	fmt.Println(query)
	err = mysqlClient.QueryRow(query).Scan(&DeviceToken)
	if err != nil {
		fmt.Println("Register 查询数据库单条用户信息 出错：", err)
	}

	regIds := []string{DeviceToken}

	var code int
	var payload = &Payload{
		PushTitle:    param.Title,
		PushBody:     param.Desc,
		IsShowNotify: strconv.Itoa(param.IsShowNotify),
		Ext:          "",
	}
	switch channel {
	case "mi":
		code = miGroupPush(regIds, payload, appId)
		break
	case "hw":
		code = hwPush(regIds, payload, appId)
		break
	case "ios":
		code = iOSPush(DeviceToken, payload, appId)
		break
	case "mz":
		code = mzPush(regIds, payload, appId)
		break
	}
	if code == 1 {
		// 定义接口返回data
		data := make(map[string]interface{})
		controllers.Response(c, common.HTTPOK, "发送成功！", data)
	}
}

//消息payload，根据业务自定义
type Payload struct {
	PushTitle    string `json:"push_title"`
	PushBody     string `json:"push_body"`
	IsShowNotify string `json:"is_show_notify"`
	Ext          string `json:"ext"`
}

func miGroupPush(regIds []string, payload *Payload, appId string) int {
	appSecret := viper.GetString("mi.appSecret")
	restrictedPackageName := viper.GetString("mi.restrictedPackageName")
	payloadStr, _ := json.Marshal(payload)
	//是否透传
	passThrough := "1"
	if payload.IsShowNotify == "1" {
		passThrough = "0" //通知栏消息
	}

	var message = XMPushSDK.InitMessage(payload.PushTitle, payload.PushBody, restrictedPackageName, string(payloadStr), passThrough)

	result, err := XMPushSDK.SendRegIds(appSecret, message, regIds)
	fmt.Println(result)
	if err != nil {
		fmt.Println("群发推送报错：", err)
	}
	if result != nil && result.Code != XMPushSDK.Success {
		fmt.Println("群发推送失败，失败原因：", result.Description)
		return 0
	}
	return 1
}

func miSinglePush(regId string, payload *Payload, appId string) int {
	appSecret := viper.GetString("mi.appSecret")
	restrictedPackageName := viper.GetString("mi.restrictedPackageName")
	payloadStr, _ := json.Marshal(payload)
	//是否透传
	passThrough := "1"
	if payload.IsShowNotify == "1" {
		passThrough = "0" //通知栏消息
	}
	var message = XMPushSDK.InitMessage(payload.PushTitle, payload.PushBody, restrictedPackageName, string(payloadStr), passThrough)

	result, err := XMPushSDK.SendToOneRegId(appSecret, message, regId)
	fmt.Println(result)
	if err != nil {
		fmt.Println("群发推送报错：", err)
	}
	if result != nil && result.Code != XMPushSDK.Success {
		fmt.Println("群发推送失败，失败原因：", result.Description)
		return 0
	}
	return 1
}

func hwPush(tokens []string, payload *Payload, appId string) int {
	clientSecret := viper.GetString("hw.clientSecret")
	//restrictedPackageName := viper.GetString("hw.restrictedPackageName")
	//payloadStr, _ := json.Marshal(payload)
	//是否透传
	passThrough := "1"
	if payload.IsShowNotify == "1" {
		passThrough = "0" //通知栏消息
	}
	var message = HWPushSDK.InitMessage(payload.PushTitle, payload.PushBody, passThrough, tokens)
	fmt.Println("message:----", message)
	result, err := HWPushSDK.MessagesSend(message, appId, clientSecret)
	fmt.Println(result)
	if err != nil {
		fmt.Println("群发推送报错：", err)
	}
	if result != nil && result.Code != HWPushSDK.Success {
		fmt.Println("群发推送失败，失败原因：", result.Msg)
		return 0
	}
	return 1
}

func iOSPush(regId string, payload *Payload, appId string) int {
	keyId := viper.GetString("ios.keyId")
	teamId := viper.GetString("ios.teamId")
	bundleID := viper.GetString("ios.bundleID")
	fmt.Println(keyId, teamId)
	//是否透传
	passThrough := "1"
	if payload.IsShowNotify == "1" {
		passThrough = "0" //通知栏消息
	}
	var message = iOSPushSDK.InitMessage(payload.PushTitle, payload.PushBody, passThrough)
	authToken, err := iOSPushSDK.GetAuthToken("./config/iosP8/AuthKey_R66FMTN5B2.p8", keyId, teamId)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("aaaaaa:", authToken)
	result, err := iOSPushSDK.MessagesSend(message, regId, authToken, bundleID)
	fmt.Println(result)
	if err != nil {
		fmt.Println("群发推送报错：", err)
	}
	if result != nil && result.Status != iOSPushSDK.Success {
		fmt.Println("群发推送失败，失败原因：", result.Reason)
		return 0
	}
	return 1
}

func mzPush(regId []string, payload *Payload, appId string) int {
	appSecret := viper.GetString("mz.appSecret")
	fmt.Println(appId, appSecret)

	var message = MZPushSDK.InitMessage(payload.PushTitle, payload.PushBody)
	result, err := MZPushSDK.MessageSend(message, appId, regId, appSecret)
	fmt.Println(result)
	if err != nil {
		fmt.Println("群发推送报错：", err)
	}
	if result != nil && result.Code != MZPushSDK.Success {
		fmt.Println("群发推送失败，失败原因：", result.Message)
		return 0
	}
	return 1
}

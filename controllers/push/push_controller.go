package push

import (
	"InformationPush/common"
	"InformationPush/controllers"
	"InformationPush/lib/mysqllib"
	"InformationPush/lib/redislib"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
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

	// 根据appid获取平台参数（1,缓存redis（hash结构）、2,mysql
	var plat PushSDK.PlatformParam
	appIdKey := "AppIdKey:" + common.SMD5("AppId", appId)
	redisClient := redislib.GetClient()
	platRedis, err := redisClient.HGetAll(ctx, appIdKey).Result()
	if err != nil {
		fmt.Println("err", err)
	} else if len(platRedis) == 0 {
		fmt.Println("platform param first in redis info")
		query1 := "SELECT hw_appId, hw_clientSecret, iOS_keyId, iOS_teamId, iOS_bundleId, iOS_authTokenPath, mi_appSecret, mi_restrictedPackageName, mz_appId, mz_appSecret, oppo_appKey, oppo_masterSecret FROM platform_param WHERE app_id = '" + "1" + "'"
		err = mysqlClient.QueryRow(query1).Scan(&plat.HWAppId, &plat.HWClientSecret, &plat.IOSKeyId, &plat.IOSTeamId, &plat.IOSBundleId, &plat.IOSAuthTokenPath, &plat.MIAppSecret, &plat.MIRestrictedPackageName, &plat.MZAppId, &plat.MZAppSecret, &plat.OPPOAppKey, &plat.OPPOMasterSecret)
		if err != nil {
			fmt.Println("Platform 参数 查询数据库单条用户信息 出错：", err)
			controllers.Response(c, common.HTTPError, "no platform parma", data)
			return
		}
		// 遍历结构体，存入redis
		t := reflect.TypeOf(plat)
		v := reflect.ValueOf(plat)
		for k := 0; k < t.NumField(); k++ {
			redisClient.HSet(ctx, appIdKey, fmt.Sprint(t.Field(k).Name), fmt.Sprint(v.Field(k).Interface()))
		}
	} else {
		fmt.Println("get platform param from redis", platRedis, reflect.TypeOf(platRedis))
		// map 转 struct
		if err = mapstructure.Decode(platRedis, &plat); err != nil {
			fmt.Println(err)
		}
	}

	switch channel {
	case "hw":
		send.SetHWAppId(plat.HWAppId).SetHWClientSecret(plat.HWClientSecret)
		break
	case "ios":
		// 获取iOS-authtoken 不能频繁刷新，间隔时间为20分钟
		iOSAuthTokenKey := "iOSAuthTokenKey:" + plat.IOSKeyId + plat.IOSTeamId
		iOSAuthRedis, err := redisClient.Get(ctx, iOSAuthTokenKey).Result()
		if err == redis.Nil || iOSAuthRedis == "" {
			fmt.Println("ios-authToken first in redis info")
			authToken, err := PushSDK.GetAuthTokenIOS(plat.IOSAuthTokenPath, plat.IOSKeyId, plat.IOSTeamId)
			if err != nil {
				fmt.Println(err)
			}
			redisClient.Set(ctx, iOSAuthTokenKey, authToken, 20*time.Minute)
			send.SetIOSAuthToken(authToken)
			plat.IOSAuthToken = authToken
		} else {
			send.SetIOSAuthToken(iOSAuthRedis)
			plat.IOSAuthToken = iOSAuthRedis
		}
		send.SetIOSBundleId(plat.IOSBundleId)
		break
	case "mi":
		send.SetMIAppSecret(plat.MIAppSecret).SetMIRestrictedPackageName(plat.MIRestrictedPackageName)
		break
	case "mz":
		send.SetMZAppId(plat.MZAppId).SetMZAppSecret(plat.MZAppSecret)
		break
	case "oppo":
		send.SetOPPOAppKey(plat.OPPOAppKey).SetOPPOMasterSecret(plat.OPPOMasterSecret)
		break
	case "vivo":
		// 获取vivo-authtoken 不能频繁刷新，间隔时间为60分钟
		vAuthTokenKey := "vAuthTokenKey:" + plat.VIAppID + plat.VIAppKey
		vAuthRedis, err := redisClient.Get(ctx, vAuthTokenKey).Result()
		if err == redis.Nil || vAuthRedis == "" {
			fmt.Println("v-authToken first in redis info")
			vAuthToken := PushSDK.GetAuthTokenV(plat.VIAppID, plat.VIAppKey, plat.VIAppSecret)
			if err != nil {
				fmt.Println(err)
			}
			redisClient.Set(ctx, vAuthTokenKey, vAuthToken, 60*time.Minute)
			send.SetVIAuthToken(vAuthToken)
			plat.VIAuthToken = vAuthToken
		} else {
			send.SetVIAuthToken(vAuthRedis)
			plat.VIAuthToken = vAuthRedis
		}
		break
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

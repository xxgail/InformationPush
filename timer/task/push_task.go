package task

import (
	"InformationPush/common"
	"InformationPush/lib/redislib"
	"InformationPush/worker"
	"context"
	"fmt"
	"github.com/techoner/gophp/serialize"
	"strconv"
	"time"
)

func PushTaskInit() {
	Timer(1*time.Second, 10*time.Second, send, "", nil, nil)
}

func send(param interface{}) (result bool) {
	fmt.Println("定时器来啦！！！！！！！！！！！！")
	result = true
	key := "SendMessage"
	redisClient := redislib.GetClient()
	messageLen, err := redisClient.LLen(common.Ctx, key).Result()
	if err != nil {
		fmt.Println("获取需要发送的信息失败", err)
	}
	if messageLen == 0 {
		fmt.Println("没有需要发送的信息数据")
	}
	for messageLen > 0 {
		val, _ := redisClient.LPop(context.Background(), key).Result()
		if val == "" { // 取出来的val格式为string
			fmt.Println("val取出来的为空！")
			messageLen--
			continue
		}
		str, _ := serialize.UnMarshal([]byte(val))
		if str == nil {
			fmt.Println("str取出来的为空！")
			messageLen--
			continue
		}
		m := str.(map[string]interface{})
		param := map[string]interface{}{
			"title":   m["title"],
			"content": m["content"],
			"channel": m["channel"],
			"pushId":  m["pushId"],
			"plat":    m["plat"],
		}
		var sendTimeInt64 int64
		if m["send_time"] != nil {
			sendTime := m["send_time"].(string)
			sendTimeInt64, _ = strconv.ParseInt(sendTime, 10, 64)
		}
		worker.PushJob(param, 1, sendTimeInt64)
		messageLen--
	}
	return
}

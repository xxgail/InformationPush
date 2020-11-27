package task

import (
	"InformationPush/common"
	"InformationPush/lib/redislib"
	"InformationPush/worker"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func PushTaskInit() {
	Timer(1*time.Second, 5*time.Second, sendGroup, "", nil, nil)
}

func sendGroup(param interface{}) (result bool) {
	result = true
	key := "SendMessage"
	redisClient := redislib.GetClient()
	messageLen, err := redisClient.LLen(common.Ctx, key).Result()
	if err != nil {
		fmt.Println("获取需要发送的信息失败", err)
	}
	if messageLen == 0 {
		fmt.Println("没有需要发送的信息数据")
		return
	}
	var params []worker.PushJobParam
	for messageLen > 0 {
		val, _ := redisClient.LPop(context.Background(), key).Result()
		if val == "" { // 取出来的val格式为string
			fmt.Println("~~~~~~~~~~~~val取出来的为空！~~~~~~~~~~~~")
			messageLen--
			continue
		}
		var m worker.PushJobParam
		_ = json.Unmarshal([]byte(val), &m)

		params = append(params, m)
		messageLen--
	}
	worker.PushJobGroup(params, 1)
	return
}

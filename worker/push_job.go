package worker

import (
	"InformationPush/common"
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"reflect"
	"strconv"
	"time"
)

func PushJob(param map[string]interface{}, retryCount int, delayTime int64) {
	var args []tasks.Arg
	var arg tasks.Arg
	for _, v := range param {
		arg = tasks.Arg{
			Type:  reflect.TypeOf(v).Name(),
			Value: v,
		}
		args = append(args, arg)
	}
	signature := workTaskInit(sendMessage, args)

	fmt.Println(common.GetFileLineNum())

	// 重复次数
	signature.RetryCount = retryCount
	eta := time.Now().UTC().Add(time.Second * time.Duration(delayTime))
	signature.ETA = &eta

	singleTask(signature)
}

type PushJobParam struct {
	SendMessageParam string `json:"send_message_param"`
	DelayTime        string `json:"delay_time"`
}

func PushJobGroup(params []PushJobParam, retryCount int) {
	var signatures []*tasks.Signature
	for _, param := range params {
		var args []tasks.Arg
		args = []tasks.Arg{
			{
				Type:  "string",
				Value: param.SendMessageParam,
			},
		}
		signature := workTaskInit(sendMessage, args)
		signature.RetryCount = retryCount

		delayTime := param.DelayTime
		delayTimeInt64, _ := strconv.ParseInt(delayTime, 10, 64)
		eta := time.Now().UTC().Add(time.Second * time.Duration(delayTimeInt64))
		signature.ETA = &eta

		signatures = append(signatures, signature)
	}
	groupTask(signatures)
}

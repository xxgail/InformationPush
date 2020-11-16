package worker

import (
	"InformationPush/common"
	"encoding/json"
	"fmt"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/xxgail/PushSDK"
	"strings"
	"time"
)

// 定义task
func SayHello(str string) (string, error) {
	fmt.Println("Helloooo", str)
	log.INFO.Println("Helloooo", str)
	return "Helloooo" + str, nil
}

func LongRunningTask() error {
	log.INFO.Print("Long running tasks started")
	for i := 0; i < 10; i++ {
		log.INFO.Print(10 - i)
		time.Sleep(1 * time.Second)
	}
	log.INFO.Print("Long running tasks finished")
	return nil
}

func SendMessage(title string, content string, channel string, pushId string, plat string) (string, error) {
	fmt.Println("我执行到这里了！！！！！！！！！")
	send := PushSDK.NewSend()
	send.SetTitle(title).SetContent(content).SetApnsId(common.GetRandomApnsId())
	send.SetPushId(strings.Split(pushId, ","))
	send.SetChannel(channel)
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
		return "", send.Err
	}
	response, err := send.SendMessage()
	responseStr, _ := json.Marshal(response)
	fmt.Println(string(responseStr))
	return string(responseStr), err
}

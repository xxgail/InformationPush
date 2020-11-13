package worker

import (
	"encoding/json"
	"fmt"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/xxgail/PushSDK"
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

func SendMessage(s string) (string, error) {
	fmt.Println("我执行到这里了！！！！！！！！！")
	m := make(map[string]string)
	_ = json.Unmarshal([]byte(s), &m)
	fmt.Println("mmmmmmmmmmmm", m)
	send := PushSDK.NewSend()
	send.SetTitle(m["title"]).SetContent(m["content"]).SetApnsId(m["apnsId"])
	send.SetPushId([]string{m["pushId"]})
	send.SetChannel(m["channel"])
	plat := m["plat"]
	switch m["channel"] {
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

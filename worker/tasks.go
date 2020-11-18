package worker

import (
	"InformationPush/common"
	"encoding/json"
	"fmt"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/xxgail/PushSDK"
	"strings"
)

// 定义task
func SayHelloTask(str string) (string, error) {
	fmt.Println("Hello", str)
	log.INFO.Println("Hello", str)
	return "Hello" + str, nil
}

type MessageTaskStruct struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Channel string `json:"channel"`
	PushId  string `json:"pushId"`
	Plat    string `json:"plat"`
}

func SendMessageTask(str string) (string, error) {
	fmt.Println("我执行到这里了！！！！！！！！！")

	var m MessageTaskStruct
	_ = json.Unmarshal([]byte(str), &m)

	message := PushSDK.NewMessage()
	message.SetTitle(m.Title).SetContent(m.Content).SetApnsId(common.GetRandomApnsId())
	if message.Err != nil {
		return "", message.Err
	}

	send := PushSDK.NewSend()
	send.SetPushId(strings.Split(m.PushId, ","))
	send.SetChannel(m.Channel)
	send.SetPlatForm(m.Plat)

	if send.Err != nil {
		return "", send.Err
	}
	response, err := send.SendMessage(message)

	responseStr, _ := json.Marshal(response)
	return string(responseStr), err
}

package main

import (
	"LogCollectAgent/kafka"
	"LogCollectAgent/model"
	"LogCollectAgent/tailf"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
)

func Run() (err error) {
	fmt.Printf("Run()...\n")
	for {
		textMsg := tailf.FetchOneMsgFromChannel()
		err = sendToKafka(textMsg)
		time.Sleep(time.Millisecond * 500)
		if err != nil {
			logs.Error("Send to kafka failed, [error=%v]", err)
			continue
		}
	}
}

func sendToKafka(textMsg *model.TextMsg) (err error) {
	fmt.Printf("Prepare to send textMsg[Msg=%v,Topic=%v]\n", textMsg.Msg, textMsg.Topic)
	err = kafka.SendMessageToKafka(textMsg.Msg, textMsg.Topic)
	return
}

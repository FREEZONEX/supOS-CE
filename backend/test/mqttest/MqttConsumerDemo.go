package main

import (
	"backend/internal/repo/event/subDev"
	"context"
	"fmt"
	"log"
	"time"

	"gitee.com/unitedrhino/share/conf"
)

type msgConsumer struct {
}

func (mc *msgConsumer) OnMsg(ctx context.Context, topic string, msgId int, payload []byte) {
	cli := ctx.Value("client")
	cliId := ""
	if cli != nil {
		if client, ok := cli.(string); ok {
			cliId = client
		} else {
			cliId = fmt.Sprintf("%v", cli)
		}
	}
	log.Printf("OnMsg[%v - cli: %v]%d: %v\n", topic, cliId, msgId, string(payload))
}
func main() {
	mc := &msgConsumer{}
	// 本地先启动mqtt broker: gmqtt.exe start
	conf := &conf.MqttConf{ClientID: "supos_test", Brokers: []string{"192.168.236.101:1883"}, ConnNum: 3}
	cli, er := subDev.NewMqttClient(conf, mc)
	if er != nil {
		log.Fatalf("NewMqttClient(%v) failed", er)
	}
	_ = cli.SubscribeAll()
	time.Sleep(1 * time.Hour)
}

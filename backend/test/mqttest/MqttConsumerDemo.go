package main

import (
	"backend/internal/repo/event/subDev"
	"context"
	"log"
	"time"

	"gitee.com/unitedrhino/share/conf"
)

type msgConsumer struct {
}

func (mc *msgConsumer) OnMsg(ctx context.Context, topic string, msgId int, payload []byte) {
	log.Printf("OnMsg[%v]%d: %v\n", topic, msgId, string(payload))
}
func main() {
	mc := &msgConsumer{}
	// 本地先启动mqtt broker: gmqtt.exe start
	conf := &conf.MqttConf{ClientID: "supos_test", Brokers: []string{"127.0.0.1:1883"}, ConnNum: 1}
	cli, er := subDev.NewMqttClient(conf, mc)
	if er != nil {
		log.Fatalf("NewMqttClient(%v) failed", er)
	}
	_ = cli.SubscribeAll()
	time.Sleep(1 * time.Hour)
}

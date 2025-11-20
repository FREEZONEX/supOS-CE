package main

import (
	"backend/internal/repo/event/subDev"
	"log"
	"testing"
	"time"

	"gitee.com/unitedrhino/share/conf"
)

func TestPublishMqtt(t *testing.T) {
	conf := &conf.MqttConf{ClientID: "supos_test", Brokers: []string{"127.0.0.1:1883"}, ConnNum: 1}
	cli, er := subDev.NewMqttClient(conf, nil)
	if er != nil {
		log.Fatalf("NewMqttClient(%v) failed", er)
	}
	msg := `{ "tm": 12321,"wq": 123.21 }`
	er = cli.Publish("akJia_seqxajk8", 0, false, []byte(msg))
	log.Println("发送消息:", er, msg)
	time.Sleep(1 * time.Second)
}

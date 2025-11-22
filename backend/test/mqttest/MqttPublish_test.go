package main

import (
	"backend/internal/repo/event/subDev"
	"fmt"
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
	msg := fmt.Sprintf(`{ "a": %d }`, time.Now().UnixMilli())
	er = cli.Publish("文件2/文件2/文件2/状态/状态/状态/天气1", 0, false, []byte(msg))
	log.Println("发送消息:", er, msg)
	time.Sleep(1 * time.Second)
}

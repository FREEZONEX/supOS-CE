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
	msg := fmt.Sprintf(`{ "msg": "H8-world-%v" }`, time.Now().String())
	er = cli.Publish("guanxi/uns_folder_type_state___EXTRA___interface_______/Retmgxh-3", 0, false, []byte(msg))
	log.Println("发送消息:", er, msg)
	time.Sleep(1 * time.Second)
}

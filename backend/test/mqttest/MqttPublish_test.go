package main

import (
	"backend/internal/repo/event/subDev"
	"fmt"
	"log"
	"testing"
	"time"

	"gitee.com/unitedrhino/share/conf"
)

func TestMills(t *testing.T) {
	mill := 1767504591002
	utcTime := time.UnixMilli(int64(mill)).UTC()
	v := utcTime.Format("2006-01-02 15:04:05.000") + "+00"
	t.Log(v)

	var m map[[2]int64]string = make(map[[2]int64]string, 8)
	m[[2]int64{1, 2}] = "12"
	m[[2]int64{2, 1}] = "21"
	m[[2]int64{1, 2}] = "120"
	m[[2]int64{2, 1}] = "210"
	t.Log(m)
}
func TestPublishMqtt(t *testing.T) {
	conf := &conf.MqttConf{ClientID: "go_test", Brokers: []string{"localhost:1883"}, ConnNum: 1}
	cli, er := subDev.NewMqttClient(conf, nil)
	if er != nil {
		log.Fatalf("NewMqttClient(%v) failed", er)
	}
	ts := time.Now().UnixMilli() //`2025-12-26T05:00:28.002Z`
	var msgs []string
	one := false
	if one {
		msg := fmt.Sprintf(`[{"timeStamp":"%v","asd":3}]`, ts)
		msgs = append(msgs, msg)

		msgs = append(msgs, fmt.Sprintf(`[{"timeStamp":"%v","ax":5}]`, ts))
	} else {
		msg := fmt.Sprintf(`[{"timeStamp":"%v","asd":30},{"timeStamp":"%v","ax":71}]`, ts, ts)
		msgs = append(msgs, msg)
	}
	for _, msg := range msgs {
		er = cli.Publish("_cunchuzhuangtai_269d188d35fd429e95da", 0, false, []byte(msg))
		t.Log("发送消息:", er, msg)
		time.Sleep(1 * time.Millisecond)
	}
}

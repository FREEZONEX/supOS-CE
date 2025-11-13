package subDev

import (
	"backend/share/clients"

	"gitee.com/unitedrhino/share/conf"
)

type (
	MqttClient struct {
		client *clients.MqttClient
	}
	//登录登出消息
	ConnectMsg struct {
		UserName string `json:"username"`
		Ts       int64  `json:"ts"`
		Address  string `json:"ipaddress"`
		ClientID string `json:"clientid"`
		Reason   string `json:"reason"`
	}
)

const (
	// ShareSubTopicPrefix emqx 共享订阅前缀 参考: https://docs.emqx.com/zh/enterprise/v4.4/advanced/shared-subscriptions.html
	ShareSubTopicPrefix = "$share/uns.rpc/"
	// TopicConnectStatus emqx 客户端上下线通知 参考: https://docs.emqx.com/zh/enterprise/v4.4/advanced/system-topic.html#客户端上下线事件
	TopicConnectStatus = ShareSubTopicPrefix + "$SYS/brokers/+/clients/#"

	TopicUns = ShareSubTopicPrefix + "#"
)

func NewMqttClient(conf *conf.MqttConf, consumer clients.MqttMsgConsumer) (*MqttClient, error) {
	if conf.ConnNum == 0 {
		conf.ConnNum = 1
	}
	mc, err := clients.NewMqttClient(conf, consumer)
	if err != nil {
		return nil, err
	}
	return &MqttClient{
		client: mc,
	}, nil
}
func (d *MqttClient) SubscribeAll() error {
	//暂不用共享的 TopicUns，实测gmqtt重连后,新消息一直会收到两份
	return d.client.Subscribe("#", 1)
}
func (d *MqttClient) Subscribe(topics map[string]byte) error {
	return d.client.SubscribeMultiple(topics)
}
func (d *MqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	return d.client.Publish(topic, qos, retained, payload)
}

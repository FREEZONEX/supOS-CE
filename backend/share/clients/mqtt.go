package clients

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"

	"gitee.com/unitedrhino/share/errors"
	"github.com/google/uuid"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
)

type MqttClient struct {
	clients         []mqtt.Client
	cfg             *conf.MqttConf
	subscribeTopics map[string]byte
	consumer        MqttMsgConsumer
	lock            sync.RWMutex
}
type MqttMsgConsumer interface {
	OnMsg(ctx context.Context, topic string, msgId int, message []byte)
}

func NewMqttClient(conf *conf.MqttConf, consumer MqttMsgConsumer) (cli *MqttClient, err error) {
	var clients []mqtt.Client
	var start = time.Now()
	for len(clients) < conf.ConnNum {
		var (
			mc mqtt.Client
		)
		cli = &MqttClient{cfg: conf, consumer: consumer, subscribeTopics: make(map[string]byte, 8)}
		var tryTime = 3
		for i := 1; i <= tryTime; i++ {
			mc, err = initMqtt(conf, cli)
			logx.Infof("mqtt_client initMqtt2 mc:%v err:%v", mc, err)
			if err != nil { //出现并发情况的时候可能联犀的http还没启动完毕
				logx.Errorf("mqtt_client 连接失败 重试剩余次数:%v", tryTime-i)
				time.Sleep(time.Second * time.Duration(i))
				continue
			}
			break
		}
		if err != nil {
			logx.Errorf("mqtt_client 连接失败 conf:%#v  err:%v", conf, err)
			//os.Exit(-1)
			return
		}
		clients = append(clients, mc)
		cli.clients = clients
	}
	logx.Infof("mqtt_client 连接完成 clientNum:%v use:%s", len(clients), time.Now().Sub(start))
	return
}

func (m *MqttClient) Subscribe(topic string, qos byte) error {
	if m.consumer == nil {
		return fmt.Errorf("consumer is nil")
	}

	m.lock.Lock()
	if _, has := m.subscribeTopics[topic]; has {
		m.lock.Unlock()
		return nil
	}
	m.subscribeTopics[topic] = qos
	m.lock.Unlock()

	logx.Infof("mqtt_client_subscribe topic:%v", topic)
	var cli = m.clients[0]
	err := cli.Subscribe(topic, qos, m.subscribeHandler).Error()
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}
func (m *MqttClient) SubscribeMultiple(filters map[string]byte) error {
	if m.consumer == nil {
		return fmt.Errorf("consumer is nil")
	}
	newAdds := make(map[string]byte, len(filters))
	m.lock.Lock()
	for t, q := range filters {
		if _, has := m.subscribeTopics[t]; !has {
			m.subscribeTopics[t] = q
		} else {
			newAdds[t] = q
		}
	}
	m.lock.Unlock()

	if len(newAdds) == 0 {
		return nil
	}
	logx.Infof("mqtt_client_subscribe topics: %+v", filters)
	var cli = m.clients[0]
	err := cli.SubscribeMultiple(newAdds, m.subscribeHandler).Error()
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}
func (m *MqttClient) subscribeHandler(client mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	if strings.HasPrefix(topic, "$") {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	utils.Recover(ctx)
	m.consumer.OnMsg(ctx, topic, int(message.MessageID()), message.Payload())
}

func (m *MqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	id := rand.Intn(len(m.clients))
	return m.clients[id].Publish(topic, qos, retained, payload).Error()
}

func initMqtt(conf *conf.MqttConf, mcs *MqttClient) (mc mqtt.Client, err error) {
	opts := mqtt.NewClientOptions()
	for _, broker := range conf.Brokers {
		opts.AddBroker(broker)
	}
	randId := uuid.NewString()
	clientID := conf.ClientID + "_" + randId
	logx.Infof("mqtt_client initMqtt conf:%#v clientID:%v brokers:%#v stack=%s", conf, clientID, opts.Servers, utils.Stack(1, 10))
	opts.SetClientID(clientID).SetUsername(conf.User).SetPassword(conf.Pass)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		er := reSubscribe(client, mcs)
		logx.Infof("mqtt_client Connected clientID:%v, er=%v", clientID, er)
	})
	opts.SetReconnectingHandler(func(client mqtt.Client, options *mqtt.ClientOptions) {
		logx.Infof("mqtt_client Reconnecting clientID:%#v", options)
	})

	opts.SetAutoReconnect(true).SetMaxReconnectInterval(30 * time.Second) //意外离线的重连参数
	//opts.SetConnectRetry(true).SetConnectRetryInterval(5 * time.Second)   //首次连接的重连参数
	opts.SetConnectRetry(false)

	opts.SetConnectionAttemptHandler(func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		logx.Infof("mqtt_client 正在尝试连接 broker:%v clientID:%v", utils.Fmt(broker), clientID)
		return tlsCfg
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logx.Errorf("mqtt_client 连接丢失 err:%v  clientID:%v", utils.Fmt(err), clientID)
	})
	mc = mqtt.NewClient(opts)
	er2 := mc.Connect().WaitTimeout(5 * time.Second)
	if er2 == false || !mc.IsConnected() {
		logx.Errorf("mqtt_client 连接失败超时")
		err = fmt.Errorf("mqtt_client 连接失败")
		return nil, err
	}
	return
}

func reSubscribe(client mqtt.Client, mcs *MqttClient) (err error) {
	if mcs.consumer != nil {
		mcs.lock.RLock()
		err = client.SubscribeMultiple(mcs.subscribeTopics, mcs.subscribeHandler).Error()
		mcs.lock.RUnlock()
	}
	return err
}

package common

import (
	"backend/internal/repo/event/subDev"
	"backend/internal/svc"
	"fmt"
	"sync"
)

var (
	unsMqttPublisher   *subDev.MqttClient
	unsMqttPublisherMu sync.Mutex
)

// PublishUnsMessage publishes a payload to the configured UNS MQTT broker.
func PublishUnsMessage(svcCtx *svc.ServiceContext, topic string, payload []byte) error {
	if len(topic) == 0 {
		return fmt.Errorf("mqtt publish topic is empty")
	}
	publisher, err := getUnsMqttPublisher(svcCtx)
	if err != nil {
		return err
	}
	return publisher.Publish(topic, 1, false, payload)
}

func getUnsMqttPublisher(svcCtx *svc.ServiceContext) (*subDev.MqttClient, error) {
	if svcCtx == nil {
		return nil, fmt.Errorf("service context is nil")
	}
	if svcCtx.Config.DevLink.Mode != "mqtt" {
		return nil, fmt.Errorf("mqtt publish unavailable: devlink.mode=%s", svcCtx.Config.DevLink.Mode)
	}
	if len(svcCtx.Config.DevLink.Mqtt.Brokers) == 0 {
		return nil, fmt.Errorf("mqtt publish unavailable: no brokers configured")
	}

	unsMqttPublisherMu.Lock()
	defer unsMqttPublisherMu.Unlock()

	if unsMqttPublisher != nil {
		return unsMqttPublisher, nil
	}

	publisher, err := subDev.NewMqttClient(&svcCtx.Config.DevLink.Mqtt, nil)
	if err != nil {
		return nil, err
	}
	unsMqttPublisher = publisher
	return unsMqttPublisher, nil
}

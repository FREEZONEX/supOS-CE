package mount

// MountMetaQueryType represents mount metadata query types
type MountMetaQueryType string

const (
	// Collector related
	MountMetaCollectorVersion MountMetaQueryType = "COLLECTOR_VERSION" // 采集器版本
	MountMetaCollector        MountMetaQueryType = "COLLECTOR"         // 采集器
	MountMetaCollectorDevice  MountMetaQueryType = "COLLECTOR_DEVICE"  // 采集器设备
	MountMetaCollectorTag     MountMetaQueryType = "COLLECTOR_TAG"     // 测点位号

	// MQTT related
	MountMetaMQTTBroker MountMetaQueryType = "MQTT_BROKER" // MQTT连接信息
	MountMetaMQTTTopic  MountMetaQueryType = "MQTT_TOPIC"  // MQTT topic

	// Kafka related
	MountMetaKafkaBroker MountMetaQueryType = "KAFKA_BROKER" // Kafka连接信息
	MountMetaKafkaTopic  MountMetaQueryType = "KAFKA_TOPIC"  // Kafka topic

	// RabbitMQ related
	MountMetaRabbitMQBroker MountMetaQueryType = "RABBITMQ_BROKER" // RabbitMQ连接信息
	MountMetaRabbitMQTopic  MountMetaQueryType = "RABBITMQ_TOPIC"  // RabbitMQ topic
)

// String returns the string representation
func (m MountMetaQueryType) String() string {
	return string(m)
}

package mount

// MountSubSourceType represents mount sub-source types
type MountSubSourceType string

const (
	MountSubSourceCollector       MountSubSourceType = "collector"        // 挂载采集器源
	MountSubSourceCollectorDevice MountSubSourceType = "collector_device" // 挂载采集器设备源
	MountSubSourceMQTTAll         MountSubSourceType = "mqtt_all"         // MQTT全挂载
	MountSubSourceKafkaAll        MountSubSourceType = "kafka_all"        // Kafka全挂载
	MountSubSourceRabbitMQAll     MountSubSourceType = "rabbitmq_all"     // RabbitMQ全挂载
)

// Type returns the type string
func (m MountSubSourceType) Type() string {
	return string(m)
}

// String returns the string representation
func (m MountSubSourceType) String() string {
	return string(m)
}

// GetMountSubSourceTypeFromType returns MountSubSourceType from type string
func GetMountSubSourceTypeFromType(typeStr string) (MountSubSourceType, bool) {
	switch MountSubSourceType(typeStr) {
	case MountSubSourceCollector, MountSubSourceCollectorDevice, MountSubSourceMQTTAll,
		MountSubSourceKafkaAll, MountSubSourceRabbitMQAll:
		return MountSubSourceType(typeStr), true
	default:
		return "", false
	}
}

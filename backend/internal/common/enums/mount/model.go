package mount

// MountModel represents mount model/method types
type MountModel string

const (
	MountModelCollectorAll    MountModel = "collector_all"    // 采集器全挂载
	MountModelCollectorDevice MountModel = "collector_device" // 采集器部分设备挂载
	MountModelMQTTAll         MountModel = "mqtt_all"         // MQTT全挂载
	MountModelKafkaAll        MountModel = "kafka_all"        // Kafka全挂载
	MountModelRabbitMQAll     MountModel = "rabbitmq_all"     // RabbitMQ全挂载
)

// Type returns the type string
func (m MountModel) Type() string {
	return string(m)
}

// String returns the string representation
func (m MountModel) String() string {
	return string(m)
}

// GetMountModelFromType returns MountModel from type string
func GetMountModelFromType(typeStr string) (MountModel, bool) {
	switch MountModel(typeStr) {
	case MountModelCollectorAll, MountModelCollectorDevice, MountModelMQTTAll,
		MountModelKafkaAll, MountModelRabbitMQAll:
		return MountModel(typeStr), true
	default:
		return "", false
	}
}

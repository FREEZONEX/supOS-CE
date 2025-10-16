package mount

// MountSourceType represents mount source types
type MountSourceType string

const (
	MountSourceCollector MountSourceType = "collector"
	MountSourceMQTT      MountSourceType = "mqtt"
	MountSourceKafka     MountSourceType = "kafka"
	MountSourceRabbitMQ  MountSourceType = "rabbitmq"
	MountSourceConnect   MountSourceType = "connect"
)

type mountSourceInfo struct {
	Type      string
	TypeValue int
}

var mountSourceDetails = map[MountSourceType]mountSourceInfo{
	MountSourceCollector: {Type: MountSourceCollector.String(), TypeValue: 16},
	MountSourceMQTT:      {Type: MountSourceMQTT.String(), TypeValue: 50},
	MountSourceKafka:     {Type: MountSourceKafka.String(), TypeValue: 51},
	MountSourceRabbitMQ:  {Type: MountSourceRabbitMQ.String(), TypeValue: 52},
	MountSourceConnect:   {Type: MountSourceConnect.String(), TypeValue: 100},
}

// String returns the string representation.
func (m MountSourceType) String() string {
	return string(m)
}

// Value returns the integer representation of the mount source type.
func (m MountSourceType) Value() int {
	if info, ok := mountSourceDetails[m]; ok {
		return info.TypeValue
	}
	return -1
}

// GetMountSourceTypeFromType returns MountSourceType from type string.
func GetMountSourceTypeFromType(typeStr string) (MountSourceType, bool) {
	for m, info := range mountSourceDetails {
		if info.Type == typeStr {
			return m, true
		}
	}
	return "", false
}

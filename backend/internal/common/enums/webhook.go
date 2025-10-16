package enums

// WebhookSubscribeEvent represents webhook subscription event types
type WebhookSubscribeEvent string

const (
	WebhookEventInstanceFieldChange WebhookSubscribeEvent = "INSTANCE_FIELD_CHANGE" // 实例字段变化事件
	WebhookEventInstanceDelete      WebhookSubscribeEvent = "INSTANCE_DELETE"       // 删除实例事件
)

// String returns the string representation
func (wse WebhookSubscribeEvent) String() string {
	return string(wse)
}

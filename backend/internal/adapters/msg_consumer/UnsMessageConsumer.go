package msg_consumer

import (
	"backend/internal/common/event"
	"backend/internal/common/service"
	"backend/internal/repo/event/subDev"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"strconv"
	"unicode"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsMessageConsumer struct {
	log        logx.Logger
	defService service.IUnsDefinitionService
}

func init() {
	spring.RegisterLazy[*UnsMessageConsumer](func() *UnsMessageConsumer {
		return &UnsMessageConsumer{
			log:        logx.WithContext(context.Background()),
			defService: spring.GetBean[service.IUnsDefinitionService](),
		}
	})
}
func (u *UnsMessageConsumer) OnMsg(ctx context.Context, topic string, msgId int, payload []byte) {
	var def *types.CreateTopicDto
	if unicode.IsDigit(rune(topic[0])) {
		idLong, numErr := strconv.ParseInt(topic, 10, 64)
		if numErr == nil {
			def = u.defService.GetDefinitionById(idLong)
		}
		if def == nil {
			def = u.defService.GetDefinitionByPath(topic)
		}
	} else {
		def = u.defService.GetDefinitionByAlias(topic)
		if def == nil {
			def = u.defService.GetDefinitionByPath(topic)
		}
	}
	var uns types.CreateTopicDto
	if def != nil {
		uns = *def
	}
	u.log.Infof("OnMsg[%s]: %v, def=%+v\n", topic, string(payload), uns)
}
func (u *UnsMessageConsumer) OnBatchMessage(payloads map[string]map[string]any) {
}
func (u *UnsMessageConsumer) OnEventContextRefreshedEvent0(ev *event.ContextRefreshedEvent) {
	if sv := ev.SvcContext; sv != nil && len(sv.Config.DevLink.Mqtt.Brokers) > 0 && sv.Config.DevLink.Mode == "mqtt" {
		go func() {
			cli, er := subDev.NewMqttClient(&sv.Config.DevLink.Mqtt, u)
			if er != nil {
				u.log.Errorf("NewMqttClient(%v) failed", er)
			} else if cli != nil {
				_ = cli.SubscribeAll()
			}
		}()
	}
}

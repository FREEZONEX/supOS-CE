package msg_consumer

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/finddatautil"
	"backend/internal/repo/event/subDev"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsMessageConsumer struct {
	log         logx.Logger
	sink        serviceApi.IDataSinkService
	defService  serviceApi.IUnsDefinitionService
	wsSender    serviceApi.IWebsocketSender
	calcService UnsRealtimeCalcService
}

func init() {
	spring.RegisterLazy[*UnsMessageConsumer](func() *UnsMessageConsumer {
		return &UnsMessageConsumer{
			log:        logx.WithContext(context.Background()),
			defService: spring.GetBean[serviceApi.IUnsDefinitionService](),
		}
	})
}

// OnMsg 处理来自mqtt的单个消息
func (u *UnsMessageConsumer) OnMsg(ctx context.Context, topic string, msgId int, payload []byte) {
	var def *types.CreateTopicDto
	if unicode.IsDigit(rune(topic[0])) {
		idLong, numErr := strconv.ParseInt(topic, 10, 64)
		if numErr == nil {
			def = u.defService.GetDefinitionById(idLong)
		}
		if def == nil && (!constants.UseAliasAsTopic || strings.Contains(topic, "/")) {
			def = u.defService.GetDefinitionByPath(topic)
		}
	} else {
		if constants.UseAliasAsTopic {
			def = u.defService.GetDefinitionByAlias(topic)
			if def == nil && strings.Contains(topic, "/") {
				def = u.defService.GetDefinitionByPath(topic)
			}
		} else {
			def = u.defService.GetDefinitionByPath(topic)
			if def == nil && !strings.Contains(topic, "/") {
				def = u.defService.GetDefinitionByAlias(topic)
			}
		}
	}
	strPayload := string(payload)
	if def == nil {
		u.log.Debugf("UnknownMsg[%s]: %v\n", topic, strPayload)
		u.getWsSender().SendMessage(serviceApi.WebsocketMessage{Path: topic, Payload: strPayload})
		return
	}
	u.log.Debugf("OnMsg[%s]: %v, def=%+v\n", topic, strPayload, *def)
	var data interface{}
	err := json.Unmarshal(payload, &data)
	if err != nil {
		u.sendErrMsg(def, strPayload, err.Error())
		return
	}
	u.sendData(u.procDataAndSendWs(def, data, nil))
}

// OnMessageByAlias 处理单个消息
func (u *UnsMessageConsumer) OnMessageByAlias(alias, payload string) {
	def := u.defService.GetDefinitionByAlias(alias)
	if def == nil {
		u.log.Infof("Unknown alias[%s]: %v\n", alias, payload)
		return
	}
	var data interface{}
	bs := []byte(payload)
	err := json.Unmarshal(bs, &data)
	if err != nil {
		u.sendErrMsg(def, payload, err.Error())
		return
	}
	u.sendData(u.procDataAndSendWs(def, data, nil))
}

// OnBatchMessage 处理批量消息
func (u *UnsMessageConsumer) OnBatchMessage(payloads map[string]map[string]any) {
	messages := make([]serviceApi.TopicMessage, 0, len(payloads))
	for alias, data := range payloads {
		def := u.defService.GetDefinitionByAlias(alias)
		if def == nil {
			u.log.Debugf("Unknown alias[%s]\n", alias)
			continue
		}
		messages = u.procDataAndSendWs(def, data, messages)
	}
	u.sendData(messages)
}

func (u *UnsMessageConsumer) procDataAndSendWs(def *types.CreateTopicDto, data any, messages []serviceApi.TopicMessage) []serviceApi.TopicMessage {
	list, erMsg := procData(def, data)
	u.sendToWebsocket(def, list, "", erMsg)
	messages = append(messages, serviceApi.TopicMessage{UnsId: def.Id, DataSrcId: types.SrcJdbcType(def.DataSrcID), Data: list})

	if len(list) > 0 && len(erMsg) == 0 {
		calcDef, calcData, calcErr := u.calcService.TryCalculate(u.defService, def, list[len(list)-1])
		if calcData != nil && calcDef != nil {
			calcList := []map[string]interface{}{calcData}
			setLastData(calcList, calcDef.GetTimestampField(), calcDef.GetFieldDefines())

			u.sendToWebsocket(calcDef, calcList, "", calcErr)
			messages = append(messages, serviceApi.TopicMessage{UnsId: calcDef.Id, DataSrcId: types.SrcJdbcType(calcDef.DataSrcID), Data: calcList})
		}
	}
	return messages
}

// OnMessageByAliasOnUpdate 处理vqt消息
func (u *UnsMessageConsumer) OnMessageByAliasOnUpdate(aliasVqtMap map[string]string) {
	msgs := make([]serviceApi.TopicMessage, 0, len(aliasVqtMap))
	for alias, payload := range aliasVqtMap {
		var data interface{}
		err := json.Unmarshal([]byte(payload), &data)
		if err != nil {
			continue
		}
		def := u.defService.GetDefinitionByAlias(alias)
		if def == nil {
			u.log.Debugf("Unknown alias[%s]\n", alias)
			continue
		}
		list, erMsg := procData(def, data)
		u.sendToWebsocket(def, list, "", erMsg)
		msgs = append(msgs, serviceApi.TopicMessage{UnsId: def.Id, DataSrcId: types.SrcJdbcType(def.DataSrcID), Data: list})
	}
	u.sendData(msgs)
}
func (u *UnsMessageConsumer) sendData(unsData []serviceApi.TopicMessage) {
	if u.sink == nil {
		u.sink = spring.GetBean[serviceApi.IDataSinkService]()
	}
	u.sink.Sink(unsData)
}

func (u *UnsMessageConsumer) sendErrMsg(def *types.CreateTopicDto, payload string, errMsg string) {
	u.sendToWebsocket(def, nil, payload, errMsg)
}
func (u *UnsMessageConsumer) sendToWebsocket(def *types.CreateTopicDto, data []map[string]any, payload string, errMsg string) {
	var lastData map[string]any
	if len(data) > 0 {
		lastData = data[0]
	}
	u.getWsSender().SendMessage(serviceApi.WebsocketMessage{Def: def, Data: lastData, Payload: payload, ErrMsg: errMsg})
}
func (u *UnsMessageConsumer) getWsSender() serviceApi.IWebsocketSender {
	if u.wsSender == nil {
		u.wsSender = spring.GetBean[serviceApi.IWebsocketSender]()
	}
	return u.wsSender
}
func procData(def *types.CreateTopicDto, data any) (list []map[string]interface{}, errMsg string) {
	fds := def.GetFieldDefines()
	CT := def.GetTimestampField()
	rs := finddatautil.FindDataList(data, 1, fds)
	list = rs.List
	if Ef, Lf := rs.ErrorField, rs.ToLongField; len(list) == 0 || Ef != "" || Lf != "" {
		var qos int64
		fieldName := ""
		if Ef != "" {
			qos = 0x400000000000000
			fieldName = Ef
			errMsg = I18nUtils.GetMessage("uns.invalid.type", Ef)
		}
		if Lf != "" {
			qos = 0x80000000000000 //超量程（工程单位）值"
			fieldName = Lf
			errMsg = I18nUtils.GetMessage("uns.invalid.toLong", Lf)
		}
		if qos != 0 {
			qosField := def.GetQualityField()
			fd := fds.FieldsMap[fieldName]
			var objMap = make(map[string]interface{})
			if dm, is := data.(map[string]interface{}); is {
				objMap[CT] = dm[CT]
			}
			var defVal any
			if fd != nil {
				defVal = fd.GetType().DefaultValue()
			} else {
				defVal = "0"
			}
			objMap[fieldName] = defVal
			objMap[qosField] = qos
			list = []map[string]interface{}{objMap}
		}
	}
	if len(list) == 0 {
		return
	}
	setLastData(list, CT, fds)
	return
}

func setLastData(list []map[string]interface{}, CT string, fds *types.FieldDefines) {
	now := time.Now().UnixMilli()
	lastMap := list[len(list)-1]
	var lastUpdateTime = now
	if lo, has := lastMap[CT].(int64); has {
		lastUpdateTime = lo
	}
	for _, v := range list {
		if _, hasCt := v[CT]; !hasCt {
			v[CT] = now
		}
	}
	for fieldName, v := range lastMap {
		fd := fds.FieldsMap[fieldName]
		if fd != nil {
			vStr := ""
			if str, isStr := v.(string); isStr {
				vStr = str
			} else {
				vStr = fmt.Sprintf("%v", v)
			}
			fd.LastValue = vStr
			fd.LastTime = lastUpdateTime
		}
	}
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

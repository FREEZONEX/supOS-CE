package msg_consumer

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/datetimeutils"
	"backend/internal/common/utils/finddatautil"
	"backend/internal/repo/event/subDev"
	"backend/internal/types"
	"backend/share/base"
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
	spring.RegisterBean[*UnsMessageConsumer](&UnsMessageConsumer{
		log: logx.WithContext(context.Background()),
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
	u.sendData(u.procDataAndSendWs(def, data, strPayload, nil))
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
	u.sendData(u.procDataAndSendWs(def, data, payload, nil))
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
		messages = u.procDataAndSendWs(def, data, "", messages)
	}
	u.sendData(messages)
}

func (u *UnsMessageConsumer) procDataAndSendWs(def *types.CreateTopicDto, data any, strPayload string, messages []serviceApi.TopicMessage) []serviceApi.TopicMessage {
	list, erMsg := procData(def, data)
	u.sendToWebsocket(def, list, strPayload, erMsg)
	save2db := base.P2v(def.Save2Db)
	if len(list) > 0 && save2db {
		messages = append(messages, serviceApi.TopicMessage{UnsId: def.Id, DataSrcId: types.SrcJdbcType(def.DataSrcID), Data: list})
	}

	if len(list) > 0 && len(erMsg) == 0 {
		calcDef, calcData, calcErr := u.calcService.TryCalculate(u.defService, def, list[len(list)-1])
		if calcData != nil && calcDef != nil {
			calcList := []map[string]interface{}{calcData}
			setLastData(calcList, calcDef)

			u.sendToWebsocket(calcDef, calcList, "", calcErr)
			if save2db {
				messages = append(messages, serviceApi.TopicMessage{UnsId: calcDef.Id, DataSrcId: types.SrcJdbcType(calcDef.DataSrcID), Data: calcList})
			}
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
		msgs = u.procDataAndSendWs(def, data, "", msgs)
	}
	u.sendData(msgs)
}
func (u *UnsMessageConsumer) sendData(unsData []serviceApi.TopicMessage) {
	if len(unsData) > 0 {
		if u.sink == nil {
			u.sink = spring.GetBean[serviceApi.IDataSinkService]()
		}
		u.sink.Sink(unsData)
	}
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
	if base.P2v(def.DataType) == constants.JsonbType {
		jsonbFiled := "json"
		if vm, ok := data.(map[string]any); ok {
			if _, has := vm[jsonbFiled]; !has {
				bs, _ := json.Marshal(data)
				vm = map[string]any{jsonbFiled: string(bs)}
			}
			list = []map[string]any{vm}
			list = setLastData(list, def)
			return
		}
	}
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
	list = setLastData(list, def)
	return
}

func setLastData(list []map[string]interface{}, def *types.CreateTopicDto) []map[string]interface{} {
	if len(list) == 0 {
		return list
	}
	CT, qos, fds := def.GetTimestampField(), def.GetQualityField(), def.GetFieldDefines()
	now := time.Now().UnixMilli()
	var lastUpdateTime = now
	var lastMap map[string]interface{}
	mergeTime := def.GetSrcJdbcType().TypeCode() == constants.TimeSequenceType
	if mergeTime {
		var prevBean map[string]interface{}
		for _, f := range def.Fields {
			if lv := f.LastValue; lv != nil {
				if prevBean == nil {
					prevBean = make(map[string]interface{}, 8)
				}
				prevBean[f.Name] = lv
			}
		}
		mergeList := mergeBeansWithTimestamp(list, CT, now, prevBean)
		if len(qos) > 0 {
			for _, vm := range mergeList {
				if _, hasQos := vm[qos]; !hasQos {
					vm[qos] = 0 // 写入质量码默认值：Good(0)
				}
			}
		}
		if len(mergeList) == 0 {
			logx.Errorf("合并数据出问题[ %s ]: %+v\n", def.Alias, list)
			return list
		}
		logx.Debugf("MergeList[ %s ]: %+v, list: %+v\n", def.Alias, mergeList, list)
		list = mergeList
	} else {
		for _, vm := range list {
			if _, hasCt := vm[CT]; !hasCt {
				vm[CT] = now
			}
		}
	}

	lastMap = list[len(list)-1]
	if lo, has := lastMap[CT].(int64); has {
		lastUpdateTime = lo
	}
	for fieldName, v := range lastMap {
		fd := fds.FieldsMap[fieldName]
		if fd != nil {
			fd.LastValue = v
			fd.LastTime = lastUpdateTime
		}
	}
	return list
}
func parseTimestamp(curT any) (ct int64) {
	if Float, isFloat := curT.(float64); isFloat { // json unmarshal 来的都是 float64 类型
		ct = int64(Float)
	} else if Long, isLong := curT.(int64); isLong {
		ct = Long
	} else {
		str := fmt.Sprint(curT)
		Double, err := strconv.ParseFloat(str, 64)
		if err != nil {
			ct = -1
			if dt, dtEr := datetimeutils.ParseDate(str); dtEr == nil && dt.Year() > 1970 {
				ct = dt.UnixMilli()
			}
		} else if Int := int64(Double); Int > 1100000000000 {
			ct = Int
		}
	}
	if ct < 1100000000000 || ct > 11000000000001 {
		return -1
	}
	return ct
}
func mergeBeansWithTimestamp(list []map[string]interface{}, CT string, now int64, prevBean map[string]any) []map[string]interface{} {
	prevTime := int64(-1)
	if len(prevBean) > 0 {
		prevTime = parseTimestamp(prevBean[CT])
	}
	mergeList := make([]map[string]interface{}, 0, len(list))
	for _, vm := range list {
		if curT, hasCt := vm[CT]; !hasCt {
			vm[CT] = now
			mergeList = append(mergeList, vm)
		} else {
			ct := parseTimestamp(curT)
			if ct < 1 {
				logx.Debugf("BadTimestamp: %v", curT)
				vm[CT] = now
				mergeList = append(mergeList, vm)
				continue
			}
			if sz := len(mergeList); ct == prevTime {
				if sz > 0 {
					last := mergeList[sz-1]
					mm := make(map[string]any, len(vm)*2)
					for k, v := range last {
						mm[k] = v
					}
					for k, v := range vm {
						mm[k] = v
					}
					mergeList[sz-1] = mm
				} else {
					var mm = vm
					if len(prevBean) > 0 {
						mm = make(map[string]any, len(vm)*2)
						for k, v := range prevBean {
							mm[k] = v
						}
						for k, v := range vm {
							mm[k] = v
						}
					}
					mergeList = append(mergeList, mm)
				}
			} else {
				mergeList = append(mergeList, vm)
			}
			prevTime = ct
			prevBean = vm
		}
	}
	return mergeList
}
func (u *UnsMessageConsumer) OnEventContextRefreshedEvent10(ev *event.ContextRefreshedEvent) {
	u.defService = spring.GetBean[*UnsDefinitionService]()
	if sv := ev.SvcContext; sv != nil && len(sv.Config.DevLink.Mqtt.Brokers) > 0 && sv.Config.DevLink.Mode == "mqtt" {
		go func() {
			cli, er := subDev.NewMqttClient(&sv.Config.DevLink.Mqtt, u)
			if er != nil {
				u.log.Errorf("NewMqttClient(%v) failed", er)
				for i := int64(5); ; i <<= 1 {
					if i < 0 {
						i = 60
					}
					time.Sleep(time.Duration(i) * time.Second)
					cli, er = subDev.NewMqttClient(&sv.Config.DevLink.Mqtt, u)
					if cli != nil && er == nil {
						break
					}
				}
			}
			if cli != nil {
				_ = cli.SubscribeAll()
			}
		}()
	}
}

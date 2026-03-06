package msg_consumer

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/datetimeutils"
	"backend/internal/repo/event/subDev"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unsafe"

	"gitee.com/unitedrhino/share/utils"
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
	var def *types.UnsDefinition
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
	payload = bytes.TrimLeftFunc(payload, unicode.IsSpace)
	strPayload := b2s(payload)
	if def == nil || len(payload) == 0 {
		u.log.Debugf("UnknownMsg[%s]: %v\n", topic, strPayload)
		u.getWsSender().SendMessage(serviceApi.WebsocketMessage{Path: topic, Payload: strPayload})
		return
	}
	u.log.Debugf("OnMsg[%s]: %v, def=%+v\n", topic, strPayload, *def)
	var t0, t1, t2 time.Time
	t0 = time.Now()
	data, err := parseJsonList(payload)
	if err != nil {
		u.sendErrMsg(def, strPayload, err.Error())
		return
	}
	ctx = context.WithValue(ctx, "payload", strPayload)
	msgList := u.procDataAndSendWs(ctx, def, data, strPayload, nil)
	t1 = time.Now()
	u.sendData(ctx, msgList)
	t2 = time.Now()
	if du := t2.Sub(t0); du > slowGap {
		logx.WithDuration(du).Slowf("sendWs: %v, sink: %v", t1.Sub(t0), t2.Sub(t1))
	}
}
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func s2b(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

const slowGap = 500 * time.Millisecond

// OnMessageByAlias 处理单个消息
func (u *UnsMessageConsumer) OnMessageByAlias(ctx context.Context, alias, payload string) {
	def := u.defService.GetDefinitionByAlias(alias)
	if def == nil {
		u.log.Infof("Unknown alias[%s]: %v\n", alias, payload)
		return
	}
	data, err := parseJsonList(s2b(payload))
	if err != nil {
		u.sendErrMsg(def, payload, err.Error())
		return
	}
	msgList := u.procDataAndSendWs(ctx, def, data, payload, nil)
	u.sendData(ctx, msgList)
}

// OnBatchMessage 处理批量消息
func (u *UnsMessageConsumer) OnBatchMessage(ctx context.Context, payloads map[string]map[string]string) {
	messages := make([]serviceApi.TopicMessage, 0, len(payloads))
	for alias, data := range payloads {
		def := u.defService.GetDefinitionByAlias(alias)
		if def == nil {
			u.log.Debugf("Unknown alias[%s]\n", alias)
			continue
		}
		messages = u.procDataAndSendWs(ctx, def, []map[string]string{data}, "", messages)
	}
	u.sendData(ctx, messages)
}

func (u *UnsMessageConsumer) procDataAndSendWs(ctx context.Context, def *types.UnsDefinition, data []map[string]string, strPayload string, messages []serviceApi.TopicMessage) []serviceApi.TopicMessage {
	list, erMsg := procData(ctx, def, data)
	u.sendToWebsocket(def, list, strPayload, erMsg)
	save2db := base.P2v(def.Save2Db)
	if len(list) > 0 && save2db {
		messages = append(messages, serviceApi.TopicMessage{UnsId: def.Id, DataSrcId: types.SrcJdbcType(def.DataSrcID), Data: list})
	}

	if len(list) > 0 && len(erMsg) == 0 {
		calcDef, calcData, calcErr := u.calcService.TryCalculate(u.defService, def, list[len(list)-1])
		if calcData != nil && calcDef != nil {
			calcList := []map[string]string{calcData}
			setLastData(ctx, calcList, calcDef)

			u.sendToWebsocket(calcDef, calcList, "", calcErr)
			if save2db {
				messages = append(messages, serviceApi.TopicMessage{UnsId: calcDef.Id, DataSrcId: types.SrcJdbcType(calcDef.DataSrcID), Data: calcList})
			}
		}
	}
	return messages
}

// OnMessageByAliasOnUpdate 处理vqt消息
func (u *UnsMessageConsumer) OnMessageByAliasOnUpdate(ctx context.Context, aliasVqtMap map[string]string) {
	msgs := make([]serviceApi.TopicMessage, 0, len(aliasVqtMap))
	for alias, payload := range aliasVqtMap {
		data, err := parseJsonList(s2b(payload))
		if err != nil {
			continue
		}
		def := u.defService.GetDefinitionByAlias(alias)
		if def == nil {
			u.log.Debugf("Unknown alias[%s]\n", alias)
			continue
		}
		msgs = u.procDataAndSendWs(ctx, def, data, "", msgs)
	}
	u.sendData(ctx, msgs)
}
func (u *UnsMessageConsumer) sendData(ctx context.Context, unsData []serviceApi.TopicMessage) {
	if len(unsData) > 0 {
		if u.sink == nil {
			u.sink = spring.GetBean[serviceApi.IDataSinkService]()
		}
		u.sink.Sink(ctx, unsData)
	}
}

func (u *UnsMessageConsumer) sendErrMsg(def *types.UnsDefinition, payload string, errMsg string) {
	u.sendToWebsocket(def, nil, payload, errMsg)
}
func (u *UnsMessageConsumer) sendToWebsocket(def *types.UnsDefinition, data []map[string]string, payload string, errMsg string) {
	var lastData map[string]string
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
func procData(ctx context.Context, def *types.UnsDefinition, data []map[string]string) (list []map[string]string, errMsg string) {
	if base.P2v(def.DataType) == constants.JsonbType {
		jsonbFiled := "json"
		vm := data[0]
		if val, has := vm[jsonbFiled]; !has {
			bs, _ := json.Marshal(data)
			vm = map[string]string{jsonbFiled: b2s(bs)}
		} else if vs := strings.TrimSpace(val); len(vs) == 0 || (vs[0] != '{' && vs[0] != '[') {
			return nil, I18nUtils.GetMessageWithCtx(ctx, "uns.invalid.json")
		}
		list = []map[string]string{vm}
		list = setLastData(ctx, list, def)
		return
	}
	errMsg = filterMsgByUns(def, data)
	if len(data) == 0 {
		return
	}
	list = setLastData(ctx, data, def)
	return
}

func setLastData(ctx context.Context, list []map[string]string, def *types.UnsDefinition) []map[string]string {
	if len(list) == 0 {
		return list
	}
	CT, fds := def.GetTimestampField(), def.GetFieldDefines()
	now := time.Now().UnixMilli()
	var lastUpdateTime = now
	mergeTime := def.GetSrcJdbcType().TypeCode() == constants.TimeSequenceType
	if mergeTime {
		var prevBean map[string]string
		prevTime := int64(-1)
		def.Lock.RLock()
		for _, f := range def.Fields {
			lv := f.LastValue
			if lv != "" {
				if prevBean == nil {
					prevBean = make(map[string]string, 8)
				}
				prevBean[f.Name] = lv
				if f.Name == CT {
					prevTime = datetimeutils.ParseTimestamp(lv)
				}
			}
		}
		def.Lock.RUnlock()

		mergeList := mergeBeansWithTimestamp(ctx, list, CT, prevTime, now, prevBean)
		if len(mergeList) == 0 {
			logx.Errorf("合并数据出问题[ %s ]: %+v\n", def.Alias, list)
			return list
		}
		logx.Debugf("MergeList[ %s ]: %+v, list: %+v\n", def.Alias, mergeList, list)
		list = mergeList
	} else {
		nowStr := strconv.FormatInt(now, 10)
		for _, vm := range list {
			if _, hasCt := vm[CT]; !hasCt {
				vm[CT] = nowStr
			}
		}
	}

	lastMap := list[len(list)-1]
	def.Lock.Lock()
	for fieldName, v := range lastMap {
		fd := fds.FieldsMap[fieldName]
		if fd != nil {
			fd.LastValue = v
			fd.LastTime = lastUpdateTime
		}
	}
	def.Lock.Unlock()
	return list
}

func mergeBeansWithTimestamp(ctx context.Context, list []map[string]string, CT string, prevTime, nowMills int64, prevBean map[string]string) []map[string]string {
	defer func() {
		if err := recover(); err != nil {
			payload := ctx.Value("payload")
			logx.WithContext(ctx).Errorf("HandleThrow|traceID=%s|error=%#v|stack=%s| CT=%s, payload=%v", utils.TraceIdFromContext(ctx), err, utils.Stack(4, 20),
				CT, payload)
		}
	}()

	now := strconv.FormatInt(nowMills, 10)
	mergeList := make([]map[string]string, 0, len(list))
	for _, vm := range list {
		if len(vm) == 0 {
			continue
		}
		if curT, hasCt := vm[CT]; !hasCt {
			vm[CT] = now
			mergeList = append(mergeList, vm)
		} else {
			ct := datetimeutils.ParseTimestamp(curT)
			if ct < 1 {
				logx.Debugf("BadTimestamp: %v", curT)
				vm[CT] = now
				mergeList = append(mergeList, vm)
				continue
			}
			if sz := len(mergeList); ct == prevTime {
				if sz > 0 {
					last := mergeList[sz-1]
					mm := make(map[string]string, len(vm)*2)
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
						mm = make(map[string]string, len(vm)*2)
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

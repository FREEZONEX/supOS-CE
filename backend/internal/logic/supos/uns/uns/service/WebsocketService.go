package service

import (
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/serviceApi"
	"backend/internal/logic/supos/uns/topology/service"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// WebsocketService manages all WebSocket connections and subscriptions
type WebsocketService struct {
	sessions           *sync.Map // map[string]*WsSubscription (sessionId -> subscription)
	idToSessionsMap    *sync.Map // map[int64]*sync.Map (unsId -> map[sessionId]bool)
	topicToSessionsMap *sync.Map // map[string]*sync.Map (topic -> map[sessionId]bool)
	aliasToSessionsMap *sync.Map // map[string]*sync.Map (alias -> map[sessionId]subValueObj)
	tpLock             sync.RWMutex
	topologySessions   map[string]bool
	unsQueryService    *UnsQueryService
	topologyService    *service.UnsTopologyService
}

func init() {
	spring.RegisterLazy[*WebsocketService](func() *WebsocketService {
		return &WebsocketService{
			sessions:           &sync.Map{},
			idToSessionsMap:    &sync.Map{},
			topicToSessionsMap: &sync.Map{},
			aliasToSessionsMap: &sync.Map{},
			topologySessions:   make(map[string]bool, 8),
			unsQueryService:    spring.GetBean[*UnsQueryService](),
			topologyService:    spring.GetBean[*service.UnsTopologyService](),
		}
	})
}

// WsSubscription represents a WebSocket subscription
type WsSubscription struct {
	conn      io.Writer
	UnsIds    *sync.Map // map[int64]bool
	Topics    *sync.Map // map[string]bool
	AliasSet  *sync.Map // map[string]bool
	WriteLock sync.Mutex
}
type wsWriter struct {
	conn *websocket.Conn
}

func (w wsWriter) Write(b []byte) (n int, err error) {
	err = w.conn.WriteMessage(websocket.TextMessage, b)
	n = len(b)
	return
}
func newWsSubscription(w io.Writer) *WsSubscription {
	return &WsSubscription{
		conn:     w,
		UnsIds:   &sync.Map{},
		Topics:   &sync.Map{},
		AliasSet: &sync.Map{},
	}
}

func (s *WebsocketService) GetSessionCount() int {
	count := 0
	s.sessions.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

func (s *WebsocketService) AddSession(sessionId string, w io.Writer) {
	s.sessions.Store(sessionId, newWsSubscription(w))
}

// TryAddSession tries to add a session with limit check (thread-safe)
// Returns true if session was added successfully, false if limit exceeded
func (s *WebsocketService) TryAddSession(sessionId string, conn *websocket.Conn, limit int) bool {
	// Use a separate mutex for session count check to avoid race condition
	// Note: This is a simplified approach. In production, consider using atomic operations
	currentCount := 0
	s.sessions.Range(func(key, value any) bool {
		currentCount++
		return true
	})

	if currentCount >= limit {
		return false
	}

	s.AddSession(sessionId, wsWriter{conn: conn})
	return true
}
func (s *WebsocketService) HasTopologies() bool {
	s.tpLock.RLock()
	defer s.tpLock.RUnlock()
	return len(s.topologySessions) > 0
}
func (s *WebsocketService) HandleSessionConnected(sessionId string, req *url.URL) {
	queryParams, err := url.ParseQuery(req.RawQuery)
	if err != nil {
		logx.Errorf("failed to parse query params: %v", err)
		return
	}

	idStrs := queryParams["id"]
	topics := queryParams["topic"]

	// Handle file import request
	if len(idStrs) == 0 && len(topics) == 0 {
		// Check for import request
		if file := queryParams.Get("file"); file != "" {
			// TODO: Handle file import (global/i18n/uns)
			logx.Infof("file import request: %s", file)
			return
		}

		// Check for topology subscription
		if globalTopology := queryParams.Get("globalTopology"); globalTopology != "" {
			subscriptionVal, _ := s.sessions.Load(sessionId)
			if subscription, ok := subscriptionVal.(*WsSubscription); ok {
				s.tpLock.Lock()
				s.topologySessions[sessionId] = true
				s.tpLock.Unlock()

				logx.Infof("topology subscription: %s", sessionId)

				// Publish initial topology message
				s.publishTopologyMessage(subscription.conn)
			}
			return
		}
		return
	}

	// Handle ID subscriptions
	if len(idStrs) > 0 {
		subscriptionVal, ok := s.sessions.Load(sessionId)
		if !ok {
			logx.Errorf("session not found: %s", sessionId)
			return
		}
		subscription := subscriptionVal.(*WsSubscription)

		logx.Infof("subscribe: %s ids=%v", sessionId, idStrs)

		for _, idStr := range idStrs {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				logx.Errorf("invalid id: %s", idStr)
				continue
			}

			subscription.UnsIds.Store(id, true)

			// Add to idToSessionsMap
			sessionsVal, _ := s.idToSessionsMap.LoadOrStore(id, &sync.Map{})
			sessions := sessionsVal.(*sync.Map)
			sessions.Store(sessionId, true)

			// Publish initial message
			s.publishMessage(subscription.conn, id)
		}
	}

	// Handle topic subscriptions
	if len(topics) > 0 {
		subscriptionVal, ok := s.sessions.Load(sessionId)
		if !ok {
			logx.Errorf("session not found: %s", sessionId)
			return
		}
		subscription := subscriptionVal.(*WsSubscription)

		logx.Infof("subscribe: %s topics=%v", sessionId, topics)

		for _, topic := range topics {
			decodedTopic, _ := url.QueryUnescape(topic)
			subscription.Topics.Store(decodedTopic, true)

			// Add to topicToSessionsMap
			sessionsVal, _ := s.topicToSessionsMap.LoadOrStore(decodedTopic, &sync.Map{})
			sessions := sessionsVal.(*sync.Map)
			sessions.Store(sessionId, true)

			// Publish initial message
			s.publishMessageByTopic(subscription.conn, decodedTopic)
		}
	}
}

func (s *WebsocketService) HandleCmdMsg(payload string, sessionId string) error {
	const SEND_PREV = "/send?t="
	const SEND_BODY = "&body="

	// Handle send command: /send?t=alias&body=payload
	if strings.HasPrefix(payload, SEND_PREV) {
		bodyIndex := strings.Index(payload[len(SEND_PREV):], SEND_BODY)
		if bodyIndex > 0 {
			alias := payload[len(SEND_PREV) : len(SEND_PREV)+bodyIndex]
			body := payload[len(SEND_PREV)+bodyIndex+len(SEND_BODY):]
			logx.Infof("ws onMessage: alias=%s, payload=%s", alias, body)
			// TODO: Call topicMessageConsumer.onMessageByAlias(alias, body)
		}
		return nil
	}

	// Handle JSON command (CMD_SUB, etc.)
	if strings.Contains(payload, "cmd") {
		var rootNode map[string]any
		if err := json.Unmarshal([]byte(payload), &rootNode); err != nil {
			s.sendCmdMessage("不是标准的json请求结构", 400, sessionId, nil)
			return err
		}

		headNode, _ := rootNode["head"].(map[string]any)
		dataNode, _ := rootNode["data"].(map[string]any)

		if headNode == nil || headNode["cmd"] == nil {
			s.sendCmdMessage("head节点不存在或cmd指令为空", 400, sessionId, rootNode)
			return nil
		}

		if dataNode == nil {
			s.sendCmdMessage("data节点不存在", 400, sessionId, rootNode)
			return nil
		}

		cmd, _ := headNode["cmd"].(float64)
		if int(cmd) == constants.CmdSub {
			// Handle subscription command
			subRealValue, ok := dataNode["sub_real_value"].(map[string]any)
			if !ok {
				s.sendCmdMessage("sub_real_value参数不存在", 400, sessionId, headNode)
				return nil
			}

			// Send subscription response
			s.sendCmdMessage("ok", 200, sessionId, headNode)

			// TODO: Push real-time data
			// s.aliasDataPush(sessionId, version, subRealValue)

			// Store subscription in aliasToSessionsMap
			subscriptionVal, ok := s.sessions.Load(sessionId)
			if ok {
				subscription := subscriptionVal.(*WsSubscription)
				for alias := range subRealValue {
					subscription.AliasSet.Store(alias, true)
					logx.Infof("subscribe: %s alias=%s", sessionId, alias)

					aliasSessionsVal, _ := s.aliasToSessionsMap.LoadOrStore(alias, &sync.Map{})
					aliasSessions := aliasSessionsVal.(*sync.Map)
					aliasSessions.Store(sessionId, subRealValue[alias])
				}
			}
		}
	}

	return nil
}

func (s *WebsocketService) HandleSessionClosed(sessionId string) {
	subscriptionVal, ok := s.sessions.Load(sessionId)
	if !ok {
		return
	}
	subscription := subscriptionVal.(*WsSubscription)

	// Remove from idToSessionsMap
	subscription.UnsIds.Range(func(key, value any) bool {
		unsId := key.(int64)
		if sessionsVal, ok := s.idToSessionsMap.Load(unsId); ok {
			sessions := sessionsVal.(*sync.Map)
			sessions.Delete(sessionId)
		}
		return true
	})

	// Remove from topicToSessionsMap
	subscription.Topics.Range(func(key, value any) bool {
		topic := key.(string)
		if sessionsVal, ok := s.topicToSessionsMap.Load(topic); ok {
			sessions := sessionsVal.(*sync.Map)
			sessions.Delete(sessionId)
		}
		return true
	})

	// Remove from aliasToSessionsMap
	subscription.AliasSet.Range(func(key, value any) bool {
		alias := key.(string)
		if sessionsVal, ok := s.aliasToSessionsMap.Load(alias); ok {
			sessions := sessionsVal.(*sync.Map)
			sessions.Delete(sessionId)
		}
		return true
	})

	// Remove from topologySessions
	s.tpLock.Lock()
	delete(s.topologySessions, sessionId)
	s.tpLock.Unlock()
	// Remove session
	s.sessions.Delete(sessionId)

	logx.Infof("session removed: %s", sessionId)
}

func (s *WebsocketService) publishMessageByTopic(conn io.Writer, topic string) {
	msg := s.getTopicLastMessageByPath(topic)
	if _, err := conn.Write(msg); err != nil {
		logx.Errorf("failed to sendWs: topic=%s, err=%v", topic, err)
	}
}

func (s *WebsocketService) publishMessage(conn io.Writer, id int64) {
	msg := s.getTopicLastMessage(id)
	if _, err := conn.Write(msg); err != nil {
		logx.Errorf("failed to sendWs: id=%d, err=%v", id, err)
	}
}

func (s *WebsocketService) publishTopologyMessage(conn io.Writer) {
	msg := s.topologyService.GetLastMsg()
	if _, err := conn.Write(msg); err != nil {
		logx.Errorf("failed to send topology message: %v", err)
	}
}

func (s *WebsocketService) sendCmdMessage(msg string, status int, sessionId string, headNode map[string]any) {
	subscriptionVal, ok := s.sessions.Load(sessionId)
	if !ok {
		return
	}
	subscription := subscriptionVal.(*WsSubscription)

	dataMap := map[string]any{
		"cmd":    constants.CmdSub,
		"msg":    msg,
		"status": status,
	}

	var response string
	if headNode != nil {
		version, _ := headNode["version"].(string)
		response = s.aliasSubResponse(version, constants.CmdSubRes, dataMap)
	} else {
		response = msg
	}

	subscription.WriteLock.Lock()
	defer subscription.WriteLock.Unlock()
	if _, err := subscription.conn.Write([]byte(response)); err != nil {
		logx.Errorf("failed to send command message: %v", err)
	}
}

func (s *WebsocketService) aliasSubResponse(version string, cmd int, dataMap map[string]any) string {
	resultJson := map[string]any{
		"head": map[string]any{
			"version": version,
			"cmd":     cmd,
		},
		"data": dataMap,
	}

	if cmd == 3 { // CMD_VAL_PUSH
		resultJson["data"] = []any{dataMap}
	}

	jsonBytes, _ := json.Marshal(resultJson)
	return string(jsonBytes)
}

// SendMessage SendLatestMsg implements the WebsocketSender interface
func (s *WebsocketService) SendMessage(wsMsg serviceApi.WebsocketMessage) {
	var unsId int64
	if def := wsMsg.Def; def != nil {
		unsId = def.Id
	}
	var path = wsMsg.Path
	if unsId != 0 {
		// Send by UNS ID
		if sessionsVal, ok := s.idToSessionsMap.Load(unsId); ok {
			sessions := sessionsVal.(*sync.Map)
			msg := processWsMsg(wsMsg)
			sessions.Range(func(key, value any) bool {
				sessionId := key.(string)
				if subscriptionVal, ok := s.sessions.Load(sessionId); ok {
					subscription := subscriptionVal.(*WsSubscription)
					subscription.WriteLock.Lock()
					defer subscription.WriteLock.Unlock()
					if _, err := subscription.conn.Write(msg); err != nil {
						logx.Errorf("fail to sendMessage to[%s], unsId=%d", sessionId, unsId)
					}
				}
				return true
			})
		}
	} else if path != "" {
		// Send by topic path
		if sessionsVal, ok := s.topicToSessionsMap.Load(path); ok {
			sessions := sessionsVal.(*sync.Map)
			msg := processWsMsg(wsMsg)

			sessions.Range(func(key, value any) bool {
				sessionId := key.(string)
				if subscriptionVal, ok := s.sessions.Load(sessionId); ok {
					subscription := subscriptionVal.(*WsSubscription)
					subscription.WriteLock.Lock()
					defer subscription.WriteLock.Unlock()
					if _, err := subscription.conn.Write(msg); err != nil {
						logx.Errorf("fail to sendMessage to[%s], topic=%s", sessionId, path)
					}
				}
				return true
			})
		}
	}
}

type TopicMessageInfo struct {
	Msg        string           `json:"msg,omitempty"`
	UpdateTime int64            `json:"updateTime,omitempty"`
	Data       map[string]any   `json:"data,omitempty"`
	Dt         map[string]int64 `json:"dt,omitempty"`
	Payload    string           `json:"payload,omitempty"`
}

func processWsMsg(message serviceApi.WebsocketMessage) []byte {
	info := &TopicMessageInfo{UpdateTime: time.Now().UnixMilli(), Msg: message.ErrMsg, Payload: message.Payload}
	if message.Def != nil {
		fs := message.Def.Fields
		if sz := len(fs); sz > 0 {
			data := make(map[string]any, sz)
			dt := make(map[string]int64, sz)
			dbType := types.SrcJdbcType(message.Def.DataSrcID).TypeCode()
			isRelation := dbType == constants.RelationType
			dm := message.Data
			hasDm := len(dm) > 0
			for _, f := range fs {
				name := f.Name
				if (isRelation && name == constants.SysFieldCreateTime) || name == constants.SystemSeqTag || name == constants.SysFieldID {
					continue
				}
				var v any
				has := false
				if hasDm {
					v, has = dm[name]
				}
				if lv := f.GetLastValue(); !has && lv != nil {
					if f.Type == types.FieldTypeDatetime {
						if date, isDate := lv.(time.Time); isDate {
							v = date.UnixMilli()
						} else if long, isLong := lv.(int64); isLong {
							v = long
						}
					} else {
						v = lv
					}
				}
				if v != nil {
					switch f.Type {
					case types.FieldTypeDouble, types.FieldTypeLong:
						v = fmt.Sprint(v)
					}
					data[name] = v
				}
				if lt := f.GetLastTime(); lt > 0 {
					dt[name] = lt
				}
			}
			if dbType == constants.TimeSequenceType {
				ct := message.Def.GetTimestampField()
				if _, has := data[ct]; !has && len(data) > 0 {
					data[ct] = info.UpdateTime
				}
			} else if base.P2v(message.Def.DataType) == constants.JsonbType {
				if jsonS, has := data["json"]; has {
					str, isStr := jsonS.(string)
					if !isStr {
						str = fmt.Sprint(jsonS)
					}
					var dMap map[string]any
					ers := json.Unmarshal([]byte(str), &dMap)
					if ers == nil {
						data = dMap
					}
				}
			}
			info.Dt = dt
			info.Data = data
		}
	}
	ret, err := json.Marshal(info)
	if err != nil {
		return emptyJson
	}
	return ret
}

var emptyJson = []byte(`{}`)

func (s *WebsocketService) getTopicLastMessage(id int64) []byte {
	if result, err := s.unsQueryService.GetLastMsg(id); err == nil && result != nil {
		return result
	}
	return emptyJson
}

func (s *WebsocketService) getTopicLastMessageByPath(topic string) []byte {
	if result, err := s.unsQueryService.GetLastMsgByPath(topic); err == nil && result != nil {
		return result
	}
	return emptyJson
}
func (s *WebsocketService) getTopicLastMessageByAlias(alias string) []byte {
	if result, err := s.unsQueryService.GetLastMsgByAlias(alias); err == nil && result != nil {
		return result
	}
	return emptyJson
}

// OnEventUnsTopologyChangeEvent handles topology change events
func (s *WebsocketService) OnEventUnsTopologyChangeEvent(e *event.UnsTopologyChangeEvent) error {
	if s.topologySessions == nil {
		return nil
	}
	// Send to all topology subscribers
	s.tpLock.RLock()
	sessionIds := base.MapKeys(s.topologySessions)
	s.tpLock.RUnlock()
	for _, sessionId := range sessionIds {
		if subscriptionVal, ok := s.sessions.Load(sessionId); ok {
			subscription := subscriptionVal.(*WsSubscription)
			subscription.WriteLock.Lock()
			if _, err := subscription.conn.Write(e.TopologyMsg); err != nil {
				logx.Errorf("fail to send topology update to session[%s]: %v", sessionId, err)
			}
			subscription.WriteLock.Unlock()
		}
	}
	return nil
}

// OnEventRemoveTopicsEvent handles topic removal events
func (s *WebsocketService) OnEventRemoveTopicsEvent(e *event.RemoveTopicsEvent) error {
	// Remove subscriptions for deleted topics
	for _, topic := range e.Topics {
		unsId := topic.GetId()
		if unsId == 0 {
			continue
		}

		// Remove from idToSessionsMap
		if sessionsVal, ok := s.idToSessionsMap.Load(unsId); ok {
			sessions := sessionsVal.(*sync.Map)
			sessions.Range(func(key, value any) bool {
				sessionId := key.(string)
				if subscriptionVal, ok := s.sessions.Load(sessionId); ok {
					subscription := subscriptionVal.(*WsSubscription)
					subscription.UnsIds.Delete(unsId)
				}
				return true
			})
			s.idToSessionsMap.Delete(unsId)
		}
	}

	logx.Infof("removed %d topic subscriptions", len(e.Topics))
	return nil
}

//// OnEventWebsocketNotifyEvent handles websocket notification events for data updates
//func (s *WebsocketService) OnEventWebsocketNotifyEvent(e *event.WebsocketNotifyEvent) error {
//	s.SendLatestMsg(e.UnsID, e.Path)
//	return nil
//}

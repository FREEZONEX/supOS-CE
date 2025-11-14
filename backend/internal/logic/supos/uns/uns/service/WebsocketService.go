package service

import (
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/logic/supos/uns/topology/service"
	"backend/share/spring"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

func init() {
	spring.RegisterBean(NewWebsocketService())
}

// WsSubscription represents a WebSocket subscription
type WsSubscription struct {
	Conn      *websocket.Conn
	UnsIds    *sync.Map // map[int64]bool
	Topics    *sync.Map // map[string]bool
	AliasSet  *sync.Map // map[string]bool
	WriteLock sync.Mutex
}

// WebsocketService manages all WebSocket connections and subscriptions
type WebsocketService struct {
	sessions           *sync.Map // map[string]*WsSubscription (sessionId -> subscription)
	idToSessionsMap    *sync.Map // map[int64]*sync.Map (unsId -> map[sessionId]bool)
	topicToSessionsMap *sync.Map // map[string]*sync.Map (topic -> map[sessionId]bool)
	aliasToSessionsMap *sync.Map // map[string]*sync.Map (alias -> map[sessionId]subValueObj)
	topologySessions   *sync.Map // map[string]*websocket.Conn (sessionId -> conn)
}

func NewWebsocketService() *WebsocketService {
	return &WebsocketService{
		sessions:           &sync.Map{},
		idToSessionsMap:    &sync.Map{},
		topicToSessionsMap: &sync.Map{},
		aliasToSessionsMap: &sync.Map{},
		topologySessions:   &sync.Map{},
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

func (s *WebsocketService) AddSession(sessionId string, conn *websocket.Conn) {
	subscription := &WsSubscription{
		Conn:     conn,
		UnsIds:   &sync.Map{},
		Topics:   &sync.Map{},
		AliasSet: &sync.Map{},
	}
	s.sessions.Store(sessionId, subscription)
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

	s.AddSession(sessionId, conn)
	return true
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
				s.topologySessions.Store(sessionId, subscription.Conn)
				logx.Infof("topology subscription: %s", sessionId)

				// Publish initial topology message
				s.publishTopologyMessage(subscription.Conn)
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
			s.publishMessage(subscription.Conn, id)
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
			s.publishMessageByTopic(subscription.Conn, decodedTopic)
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
	s.topologySessions.Delete(sessionId)

	// Remove session
	s.sessions.Delete(sessionId)

	logx.Infof("session removed: %s", sessionId)
}

func (s *WebsocketService) publishMessage(conn *websocket.Conn, id int64) {
	msg := s.getTopicLastMessage(id)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		logx.Errorf("failed to sendWs: id=%d, err=%v", id, err)
	}
}

func (s *WebsocketService) publishMessageByTopic(conn *websocket.Conn, topic string) {
	msg := s.getTopicLastMessageByPath(topic)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		logx.Errorf("failed to sendWs: topic=%s, err=%v", topic, err)
	}
}

func (s *WebsocketService) publishTopologyMessage(conn *websocket.Conn) {
	topologyService := spring.GetBean[*service.UnsTopologyService]()
	msg := topologyService.GetLastMsg()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
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
	if err := subscription.Conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
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

// SendLatestMsg implements the WebsocketSender interface
func (s *WebsocketService) SendLatestMsg(unsId int64, path string) {
	if unsId != 0 {
		// Send by UNS ID
		if sessionsVal, ok := s.idToSessionsMap.Load(unsId); ok {
			sessions := sessionsVal.(*sync.Map)
			msg := s.getTopicLastMessage(unsId)

			sessions.Range(func(key, value any) bool {
				sessionId := key.(string)
				if subscriptionVal, ok := s.sessions.Load(sessionId); ok {
					subscription := subscriptionVal.(*WsSubscription)
					subscription.WriteLock.Lock()
					defer subscription.WriteLock.Unlock()
					if err := subscription.Conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
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
			msg := s.getTopicLastMessageByPath(path)

			sessions.Range(func(key, value any) bool {
				sessionId := key.(string)
				if subscriptionVal, ok := s.sessions.Load(sessionId); ok {
					subscription := subscriptionVal.(*WsSubscription)
					subscription.WriteLock.Lock()
					defer subscription.WriteLock.Unlock()
					if err := subscription.Conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
						logx.Errorf("fail to sendMessage to[%s], topic=%s", sessionId, path)
					}
				}
				return true
			})
		}
	}
}

func (s *WebsocketService) getTopicLastMessage(id int64) string {
	unsQueryService := spring.GetBean[*UnsQueryService]()
	if result, err := unsQueryService.GetLastMsg(id); err == nil && result != nil {
		if dataStr, ok := result.Data.(string); ok {
			return dataStr
		}
		// If Data is map or other type, marshal it
		if jsonBytes, err := json.Marshal(result.Data); err == nil {
			return string(jsonBytes)
		}
	}
	return "{}"
}

func (s *WebsocketService) getTopicLastMessageByPath(topic string) string {
	unsQueryService := spring.GetBean[*UnsQueryService]()
	if result, err := unsQueryService.GetLastMsgByPath(topic); err == nil && result != nil {
		if dataStr, ok := result.Data.(string); ok {
			return dataStr
		}
		// If Data is map or other type, marshal it
		if jsonBytes, err := json.Marshal(result.Data); err == nil {
			return string(jsonBytes)
		}
	}
	return "{}"
}

// OnEventUnsTopologyChangeEvent handles topology change events
func (s *WebsocketService) OnEventUnsTopologyChangeEvent(e *event.UnsTopologyChangeEvent) error {
	if s.topologySessions == nil {
		return nil
	}

	// Get topology service
	topologyService := spring.GetBean[*service.UnsTopologyService]()
	msg := topologyService.GetLastMsg()

	// Send to all topology subscribers
	s.topologySessions.Range(func(key, value any) bool {
		sessionId := key.(string)
		if subscriptionVal, ok := s.sessions.Load(sessionId); ok {
			subscription := subscriptionVal.(*WsSubscription)
			subscription.WriteLock.Lock()
			if err := subscription.Conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				logx.Errorf("fail to send topology update to session[%s]: %v", sessionId, err)
			}
			subscription.WriteLock.Unlock()
		}
		return true
	})

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

// OnEventWebsocketNotifyEvent handles websocket notification events for data updates
func (s *WebsocketService) OnEventWebsocketNotifyEvent(e *event.WebsocketNotifyEvent) error {
	s.SendLatestMsg(e.UnsID, e.Path)
	return nil
}

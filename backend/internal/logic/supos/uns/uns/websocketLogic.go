// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"context"
	"io"
	"net/http"

	"backend/internal/common/utils/apiutil"
	"backend/internal/logic/supos/uns/uns/service"
	"backend/internal/svc"
	"backend/share/spring"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	WS_SESSION_LIMIT = 10000 // Maximum number of WebSocket sessions
)

type WebsocketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	req    *http.Request
	conn   *websocket.Conn
}

// WebSocket 连接
func NewWebsocketLogic(ctx context.Context, svcCtx *svc.ServiceContext, req *http.Request, conn *websocket.Conn) *WebsocketLogic {
	return &WebsocketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		req:    req,
		conn:   conn,
	}
}

func (l *WebsocketLogic) Websocket() {
	defer func() {
		if err := l.conn.Close(); err != nil {
			logx.Infof("failed to close websocket: %v", err)
		}
	}()

	logx.Infof("WebSocket open: %s", l.conn.RemoteAddr().String())

	// Get WebsocketService from spring container
	wsService := spring.GetBean[*service.WebsocketService]()

	// Generate unique session ID
	sessionId := uuid.New().String()

	// Check session limit with lock protection (prevent TOCTOU race condition)
	if !wsService.TryAddSession(sessionId, l.conn, WS_SESSION_LIMIT) {
		l.conn.WriteMessage(websocket.TextMessage, []byte("session reached its maximum capacity"))
		logx.Errorf("ws session exceeded limit (%d), closing connection", WS_SESSION_LIMIT)
		return
	}

	// Parse user from token (optional, if needed)
	user := apiutil.GetUserFromContext(l.ctx)
	if user != nil {
		logx.Infof("WebSocket user: %s", user.PreferredUsername)
	}

	// Handle connection established
	wsService.HandleSessionConnected(sessionId, l.req.URL)

	// Start message loop
	for {
		messageType, payload, err := l.conn.ReadMessage()
		if err != nil {
			if err == io.EOF || websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				logx.Infof("WebSocket closed normally: %s", sessionId)
			} else {
				logx.Errorf("WebSocket read error: %v", err)
			}
			break
		}

		if messageType != websocket.TextMessage {
			continue
		}

		msg := string(payload)
		logx.WithContext(l.ctx).Debugf("WebSocket handleMessage[%s]: %s", sessionId, msg)

		// Handle ping/pong heartbeat
		if msg == "ping" {
			if err := l.conn.WriteMessage(websocket.TextMessage, []byte("pong")); err != nil {
				logx.Errorf("failed to send pong: %v", err)
				break
			}
			continue
		}

		// Handle command message
		if err := wsService.HandleCmdMsg(msg, sessionId); err != nil {
			logx.Errorf("handleCmdMsg error: %v", err)
		}
	}

	// Handle connection closed
	wsService.HandleSessionClosed(sessionId)
	logx.Infof("WebSocket connection closed: %s", sessionId)
}

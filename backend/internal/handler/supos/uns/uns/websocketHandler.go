// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"net/http"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins, matching Java's setAllowedOriginPatterns("*")
	},
	ReadBufferSize:  5 * 1024 * 1024, // 5MB
	WriteBufferSize: 5 * 1024 * 1024, // 5MB
}

// WebSocket 连接
func WebsocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Errorf("websocket upgrade failed: %v", err)
			return
		}

		l := uns.NewWebsocketLogic(r.Context(), svcCtx, r, conn)
		l.Websocket()
	}
}

package uns

import (
	unsService "backend/internal/logic/supos/uns/uns/service"
	"backend/share/spring"
	"io"
	"net/url"
	"sync"

	"github.com/google/uuid"
)

var unsPushService *unsService.WebsocketService
var initOnce sync.Once

func initServiceOnce() {
	if unsPushService == nil {
		initOnce.Do(func() {
			unsPushService = spring.GetBean[*unsService.WebsocketService]()
		})
	}
}

// PushNewMsg 推送最新消息
func PushNewMsg(w io.Writer, url *url.URL) func() {
	initServiceOnce()

	rd, _ := uuid.NewRandom()
	sid := "sse_" + rd.String()
	unsPushService.AddSession(sid, w)
	go unsPushService.HandleSessionConnected(sid, url)

	return func() {
		unsPushService.HandleSessionClosed(sid)
	}
}

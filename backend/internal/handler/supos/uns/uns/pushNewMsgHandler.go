package uns

import (
	unsService "backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type sse_writer struct {
	msgChan chan string
}

func (s *sse_writer) Write(p []byte) (n int, err error) {
	s.msgChan <- string(p)
	return len(p), nil
}

// PushNewMsgHandler 推送最新消息
func PushNewMsgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		sseWriter := &sse_writer{msgChan: make(chan string, 3)}
		onClose := unsService.PushNewMsg(sseWriter, r.URL)
		keepAliveTicker := time.NewTicker(10 * time.Second)
		defer keepAliveTicker.Stop()
		// 持续监听并推送事件
		for {
			select {
			case msg := <-sseWriter.msgChan:
				// 发送事件数据
				_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
				w.(http.Flusher).Flush()
				logx.Debug("push new message to sse-writer: ", msg, err)
				if err != nil {
					cancel()
				}
			case <-keepAliveTicker.C:
				// 发送心跳/注释以保持连接
				_, err := fmt.Fprintf(w, ": keep-alive\n\n")
				if err != nil {
					cancel()
				} else {
					w.(http.Flusher).Flush()
				}
			case <-ctx.Done():
				// 客户端断开连接
				onClose()
				close(sseWriter.msgChan)
				return
			}
		}
	}
}

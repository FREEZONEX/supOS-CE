package uns

import (
	"backend/internal/common/utils/httpUtils"
	unsService "backend/internal/logic/supos/uns/uns"
	"context"
	"fmt"
	"net/http"
	"time"
	"unsafe"

	"github.com/zeromicro/go-zero/core/logx"
)

type sse_writer struct {
	closed  bool
	msgChan chan string
}

var errorSseClosed = fmt.Errorf("SseClosed")

func (s *sse_writer) Write(p []byte) (n int, err error) {
	if s.closed {
		return 0, errorSseClosed
	}
	// 非阻塞写入，如果通道满则丢弃最旧的消息
	msg := b2s(p)
	select {
	case s.msgChan <- msg:
		return len(p), nil
	default:
		// 通道满，尝试移除一个旧消息再写入
		select {
		case <-s.msgChan: // 丢弃一个旧消息
			s.msgChan <- msg // 写入新消息
		default:
			// 如果还是满，直接丢弃新消息
			logx.Error("丢弃sse消息: ", msg)
		}
		return len(p), nil
	}
}
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// PushNewMsgHandler 推送最新消息
func PushNewMsgHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("X-Accel-Buffering", "no") // 禁用Nginx缓冲

	flusher := httpUtils.HttpFlusher(w)
	_, _ = fmt.Fprint(w, "data: Connected\n\n")
	flusher.Flush()

	listen(w, r, flusher)
}

func listen(w http.ResponseWriter, r *http.Request, flusher http.Flusher) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	sseWriter := &sse_writer{msgChan: make(chan string, 100)}
	onClose := unsService.PushNewMsg(sseWriter, r.URL)
	keepAliveTicker := time.NewTicker(10 * time.Second)
	defer keepAliveTicker.Stop()
	// 持续监听并推送事件
	for {
		select {
		case msg := <-sseWriter.msgChan:
			// 发送事件数据
			_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
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
				flusher.Flush()
			}
		case <-ctx.Done():
			// 客户端断开连接
			sseWriter.closed = true
			onClose()
			close(sseWriter.msgChan)
			return
		}
	}
}

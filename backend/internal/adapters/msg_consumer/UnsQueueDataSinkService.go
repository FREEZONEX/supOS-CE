package msg_consumer

import (
	_ "backend/internal/adapters/postgresql"
	_ "backend/internal/adapters/timescaledb"
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/loggerlevel"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/diskqueue"
	"backend/share/spring"
	"context"
	"os"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsQueueDataSinkService struct {
	log                  logx.Logger
	queue                diskqueue.Interface
	run                  bool
	defService           serviceApi.IUnsDefinitionService
	once                 sync.Once
	persistentServiceMap map[types.SrcJdbcType]serviceApi.IPersistentService
}

const maxMsgSize = 4 * 1024 * 1024

func init() {
	spring.RegisterBean[*UnsQueueDataSinkService](&UnsQueueDataSinkService{
		log: logx.WithContext(context.Background()),
	})
}

func (s *UnsQueueDataSinkService) Sink(ctx context.Context, unsData []serviceApi.TopicMessage) {
	if len(unsData) > 0 {
		// 写入本地磁盘队列
		binData := encodeMsg(ctx, unsData)
		_ = s.queue.Put(binData) //TODO: 磁盘满的处理
	}
}

func (s *UnsQueueDataSinkService) OnEventShutdown(evt *event.ContextClosedEvent) {
	s.log.Infof("** UnsQueueDataSinkService.OnEventStop")
	s.run = false
	_ = s.queue.Close()
}

func (s *UnsQueueDataSinkService) OnEventStart100(evt *event.ContextRefreshedEvent) {
	s.defService = spring.GetBean[*UnsDefinitionService]()
	dir := constants.RootPath + "/queue"
	err := os.MkdirAll(dir, 666)
	if err != nil {
		panic(err)
	}
	s.queue = diskqueue.New("uns", dir,
		64*1024*1024, 8, maxMsgSize,
		2500, 5*time.Second, s.log.Debugf, s.log.Errorf)
	s.run = true
	go s.fetchData()
}

const fetchSize = 10000
const maxWait time.Duration = 1 * time.Second

func (s *UnsQueueDataSinkService) fetchData() {
	tk := time.NewTicker(maxWait)
	var size = 0
	var msgToSend = make([]serviceApi.TopicMessage, 0, fetchSize)
	for s.run {
		select {
		case <-tk.C:
			tk.Reset(maxWait)
			if size > 0 {
				//上车
				size = 0
				s.persistence(msgToSend)
				msgToSend = msgToSend[:0]
			} else if loggerlevel.DoStats {
				logx.Stat("没数据")
			}
		case msg := <-s.queue.ReadChan():
			var msgs []serviceApi.TopicMessage
			decodeMsg(msg, &msgs)
			if len(msgs) == 0 {
				continue
			}
			for _, m := range msgs {
				msgToSend = append(msgToSend, m)
				size += len(m.Data)
			}
			if size >= fetchSize {
				//上车
				size = 0
				s.persistence(msgToSend)
				msgToSend = msgToSend[:0]
			}
		}
	}
}
func (s *UnsQueueDataSinkService) persistence(msgLit []serviceApi.TopicMessage) {
	if s.persistentServiceMap == nil {
		s.once.Do(func() {
			s.persistentServiceMap = base.MapArrayToMap(spring.GetBeansOfType[serviceApi.IPersistentService](),
				func(e serviceApi.IPersistentService) (ok bool, k types.SrcJdbcType, v serviceApi.IPersistentService) {
					return true, e.GetDataSrcId(), e
				})
		})
	}
	dsMap := base.MapAndFilterGroupBy[serviceApi.TopicMessage, serviceApi.UnsData, types.SrcJdbcType](msgLit, func(e serviceApi.TopicMessage) (ok bool, id types.SrcJdbcType, dat serviceApi.UnsData) {
		def := s.defService.GetDefinitionById(e.UnsId)
		if def == nil || !base.P2v(def.Save2Db) {
			return
		}
		return true, e.DataSrcId, serviceApi.UnsData{Uns: def, Data: e.Data}
	})
	for ds, data := range dsMap {
		sv := s.persistentServiceMap[ds]
		if sv != nil {
			s.log.Debugf("Persistent[ds]: %v, len=%d", ds.Alias(), len(data))
			sv.Persistent(data)
		} else {
			s.log.Error("No persistentService: ", ds.Alias(), len(data))
		}
	}
}

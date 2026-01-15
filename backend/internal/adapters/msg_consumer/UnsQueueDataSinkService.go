package msg_consumer

import (
	_ "backend/internal/adapters/postgresql"
	_ "backend/internal/adapters/timescaledb"
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/serviceApi"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/diskqueue"
	"backend/share/spring"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
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

func (s *UnsQueueDataSinkService) Sink(unsData []serviceApi.TopicMessage) {
	if len(unsData) == 0 {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			unsDataJson, _ := json.Marshal(unsData)
			s.log.Errorf("HandleThrow|error=%#v| [%d] unsData=%v", err, len(unsDataJson), b2s(unsDataJson))
		}
	}()

	msgList := make([]*TopicMessage, len(unsData))
	for i, d := range unsData {
		dL := make([]*TopicMessage_DataArray, len(d.Data))
		for n, vm := range d.Data {
			dataMap := make(map[string]string, len(vm))
			for k, v := range vm {
				var vStr string
				if str, ok := v.(string); ok {
					vStr = str
				} else if v != nil {
					vStr = fmt.Sprint(v)
				}
				dataMap[k] = vStr
			}
			dL[n] = &TopicMessage_DataArray{Value: dataMap}
		}
		msgList[i] = &TopicMessage{Id: d.UnsId, DataSrcId: int32(d.DataSrcId), Data: dL}
	}
	list := TopicMessageList{Messages: msgList}
	binData, err := proto.Marshal(&list) //TODO: 控制 binData 不超过 maxMsgSize
	if err != nil {
		s.log.Error(err)
		return
	}
	// 写入本地磁盘队列
	_ = s.queue.Put(binData) //TODO: 磁盘满的处理

}

func (s *UnsQueueDataSinkService) OnEventShutdown(evt *event.ContextClosedEvent) {
	s.log.Infof("** UnsQueueDataSinkService.OnEventStop")
	s.run = false
	_ = s.queue.Close()
}

func (s *UnsQueueDataSinkService) OnEventStart100(evt *event.ContextRefreshedEvent) {
	s.defService = spring.GetBean[*UnsDefinitionService]()
	dir := constants.LogPath + "/queue"
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
	var msgToSend = make([]*TopicMessage, 0, fetchSize)
	for s.run {
		select {
		case <-tk.C:
			tk.Reset(maxWait)
			if size > 0 {
				//上车
				size = 0
				s.persistence(msgToSend)
				msgToSend = msgToSend[:0]
			}
		case msg := <-s.queue.ReadChan():
			var list TopicMessageList
			er := proto.Unmarshal(msg, &list)
			if er != nil {
				s.log.Error("UnmarshalErr", er)
				continue
			}
			msgs := list.Messages
			if len(msgs) == 0 {
				continue
			}
			for _, m := range msgs {
				if m != nil {
					msgToSend = append(msgToSend, m)
					size += len(m.Data)
				}
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
func (s *UnsQueueDataSinkService) persistence(msgLit []*TopicMessage) {
	if s.persistentServiceMap == nil {
		s.once.Do(func() {
			s.persistentServiceMap = base.MapArrayToMap(spring.GetBeansOfType[serviceApi.IPersistentService](),
				func(e serviceApi.IPersistentService) (ok bool, k types.SrcJdbcType, v serviceApi.IPersistentService) {
					return true, e.GetDataSrcId(), e
				})
		})
	}
	dsMap := base.MapAndFilterGroupBy[*TopicMessage, serviceApi.UnsData, types.SrcJdbcType](msgLit, func(e *TopicMessage) (ok bool, id types.SrcJdbcType, dat serviceApi.UnsData) {
		def := s.defService.GetDefinitionById(e.Id)
		if def == nil || !base.P2v(def.Save2Db) {
			return
		}
		return true, types.SrcJdbcType(e.DataSrcId), serviceApi.UnsData{Uns: def, Data: base.Map(e.Data, func(e *TopicMessage_DataArray) map[string]string {
			return e.Value
		})}
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

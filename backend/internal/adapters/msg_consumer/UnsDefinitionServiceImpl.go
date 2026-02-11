package msg_consumer

import (
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/serviceApi"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UnsDefinitionService struct {
	log                  logx.Logger
	unsMapper            dao.UnsNamespaceRepo
	cache                *ristretto.Cache[string, *types.UnsDefinition]
	locks                [1000]sync.RWMutex
	persistentServiceMap map[types.SrcJdbcType]serviceApi.IPersistentService
}

func init() {
	cache, err := ristretto.NewCache(&ristretto.Config[string, *types.UnsDefinition]{
		NumCounters: 1e6,     // number of keys to track frequency of (1M).
		MaxCost:     1 << 28, // maximum cost of cache (256M).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		logx.Errorf("init cache err: %v", err)
		panic(err)
	}
	spring.RegisterBean[*UnsDefinitionService](&UnsDefinitionService{
		log:   logx.WithContext(context.Background()),
		cache: cache,
	})
}

const keyAliasPrev = "a:"
const keyPathPrev = "p:"

func (u *UnsDefinitionService) GetDefinitionByAlias(alias string) *types.UnsDefinition {
	return u.getByAliasOrPath(keyAliasPrev, alias, u.unsMapper.GetByAlias)
}

func (u *UnsDefinitionService) GetDefinitionByPath(path string) *types.UnsDefinition {
	return u.getByAliasOrPath(keyPathPrev, path, u.unsMapper.GetByPath)
}

const costObj = 4
const costRef = 2
const costNil = 1

func (u *UnsDefinitionService) GetDefinitionById(id int64) *types.UnsDefinition {
	c := u.cache
	idStr := id2key(id)
	rs, exist := c.Get(idStr)
	if !exist {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
		defer cancel()
		db := dao.GetDb(ctx)
		po, _ := u.unsMapper.SelectById(db, id)
		if po != nil {
			rs = po2dto(po)
			c.SetWithTTL(idStr, rs, costObj, 10*time.Minute)
			key := ""
			if constants.UseAliasAsTopic {
				key = alias2key(po.Alias)
			} else {
				key = path2key(po.Path)
			}
			c.SetWithTTL(key, &types.UnsDefinition{CreateTopicDto: types.CreateTopicDto{Id: rs.Id}}, costRef, 10*time.Minute)
		} else {
			c.SetWithTTL(idStr, nil, costNil, 1*time.Minute) //占位
		}
	}
	return rs
}
func (u *UnsDefinitionService) DeleteByIds(ids []int64) error {
	for _, id := range ids {
		u.cache.Del(id2key(id))
	}
	return nil
}

func (u *UnsDefinitionService) SaveBatch(list []*types.UnsDefinition) error {
	for _, v := range list {
		u.invalidCache(v.Id, v.Alias, v.Path)
	}
	return nil
}

func (u *UnsDefinitionService) DeleteBatch(list []*types.UnsDefinition) error {
	for _, v := range list {
		u.invalidCache(v.Id, v.Alias, v.Path)
	}
	return nil
}

func (u *UnsDefinitionService) OnEventBatchCreateTableEvent0(ev *event.BatchCreateTableEvent) {
	if list := ev.Creates; len(list) > 0 {
		for _, vs := range list {
			for _, v := range vs {
				u.invalidCache(v.GetId(), v.GetAlias(), v.GetPath())
			}
		}
	}
	if list := ev.Updates; len(list) > 0 {
		for _, vs := range list {
			for _, v := range vs {
				u.invalidCache(v.GetId(), v.GetAlias(), v.GetPath())
			}
		}
	}
}
func (u *UnsDefinitionService) OnEventRemoveTopicsEvent0(ev *event.RemoveTopicsEvent) {
	if len(ev.Topics) >= 0 {
		for _, v := range ev.Topics {
			u.invalidCache(v.GetId(), v.GetAlias(), v.GetPath())
		}
	}
}
func (u *UnsDefinitionService) OnEventUpdateInstanceEvent0(ev *event.UpdateInstanceEvent) {
	if len(ev.Topics) >= 0 {
		for _, v := range ev.Topics {
			u.invalidCache(v.Id, v.Alias, v.Path)
		}
	}
}
func (u *UnsDefinitionService) invalidCache(id int64, alias, path string) {
	u.cache.Del(id2key(id))
	u.cache.Del(alias2key(alias))
	u.cache.Del(path2key(path))
}

func id2key(id int64) string {
	return strconv.FormatInt(id, 36)
}
func alias2key(alias string) string {
	return keyAliasPrev + alias
}
func path2key(path string) string {
	return keyPathPrev + path
}

func (u *UnsDefinitionService) getByAliasOrPath(kPrev string, arg string, query func(db *gorm.DB, arg string) (*dao.UnsNamespace, error)) (rs *types.UnsDefinition) {
	key := kPrev + arg
	c := u.cache

	idObj, has := c.Get(key)
	if !has {
		index := base.Abs(base.HashCode(arg)) % len(u.locks)
		u.locks[index].Lock()

		idObj, has = c.Get(key)
		if !has {
			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
			defer cancel()
			db := dao.GetDb(ctx)
			po, _ := query(db, arg)
			if po != nil {
				rs = po2dto(po)
				idKey := id2key(rs.Id)
				c.SetWithTTL(key, &types.UnsDefinition{CreateTopicDto: types.CreateTopicDto{Id: rs.Id}}, costRef, 13*time.Minute)
				c.SetWithTTL(idKey, rs, costObj, 10*time.Minute)
				c.Wait()
			} else {
				c.SetWithTTL(key, nil, costNil, 2*time.Minute) //占位
			}
			u.locks[index].Unlock()
		} else if idObj != nil {
			u.locks[index].Unlock()
			return u.GetDefinitionById(idObj.Id)
		} else {
			u.locks[index].Unlock()
		}
	} else if idObj != nil {
		return u.GetDefinitionById(idObj.Id)
	}
	return
}

func po2dto(po *dao.UnsNamespace) *types.UnsDefinition {
	rs := UnsConverter.Po2Dto(po)
	fields := rs.Fields
	def := &types.UnsDefinition{CreateTopicDto: *rs}
	if len(fields) > 0 {
		for _, field := range fields {
			field.Uns = def
		}
	}
	return def
}
func (u *UnsDefinitionService) OnEventContextRefreshedEvent1(_ *event.ContextRefreshedEvent) {
	u.persistentServiceMap = base.MapArrayToMap(spring.GetBeansOfType[serviceApi.IPersistentService](),
		func(e serviceApi.IPersistentService) (ok bool, k types.SrcJdbcType, v serviceApi.IPersistentService) {
			return true, e.GetDataSrcId(), e
		})
	types.UnsLastValueFill = u.fillUnsLastValue
}
func (u *UnsDefinitionService) fillUnsLastValue(uns *types.UnsDefinition) {
	// 查询数据库表最新的一条数据，填充字段的 lastValue
	psv, has := u.persistentServiceMap[types.SrcJdbcType(uns.DataSrcID)]
	if !has {
		return
	}
	psv.FillLastRecord(uns)
}
